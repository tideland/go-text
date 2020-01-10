// Tideland Go Text - Etc
//
// Copyright (C) 2019-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package etc // import "tideland.dev/go/text/etc"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"tideland.dev/go/dsa/collections"
	"tideland.dev/go/text/sml"
	"tideland.dev/go/text/stringex"
	"tideland.dev/go/trace/failure"
)

//--------------------
// GLOBAL
//--------------------

// key is to address a configuration inside a context.
type key int

var (
	etcKey    key = 1
	etcRoot       = []string{"etc"}
	defaulter     = stringex.NewDefaulter("etc", false)
)

//--------------------
// VALUE
//--------------------

// value helps to use the stringex.Defaulter.
type value struct {
	path    []string
	changer *collections.KeyStringValueChanger
}

// Value retrieves the value or an error. It implements
// the Valuer interface.
func (v *value) Value() (string, error) {
	sv, err := v.changer.Value()
	if err != nil {
		return "", failure.New("invalid path '%s'", fullPathToString(v.path))
	}
	return sv, nil
}

//--------------------
// ETC
//--------------------

// Application is used to apply values to a configurtation.
type Application map[string]string

// Etc contains the read etc configuration and provides access to
// it. ThetcRoot node "etc" is automatically preceded to the path.
// The node name have to consist out of 'a' to 'z', '0' to '9', and
// '-'. The nodes of a path are separated by '/'.
type Etc struct {
	values *collections.KeyStringValueTree
}

// Read reads the SML source of the configuration from a
// reader, parses it, and returns the etc instance.
func Read(source io.Reader) (*Etc, error) {
	builder := sml.NewKeyStringValueTreeBuilder()
	err := sml.ReadSML(source, builder)
	if err != nil {
		return nil, failure.Annotate(err, "invalid source format")
	}
	values, err := builder.Tree()
	if err != nil {
		return nil, failure.Annotate(err, "invalid source format")
	}
	if err = values.At("etc").Error(); err != nil {
		return nil, failure.Annotate(err, "invalid source format")
	}
	cfg := &Etc{
		values: values,
	}
	if err = cfg.postProcess(); err != nil {
		return nil, failure.Annotate(err, "cannot post-process configuration")
	}
	return cfg, nil
}

// ReadString reads the SML source of the configuration from a
// string, parses it, and returns the etc instance.
func ReadString(source string) (*Etc, error) {
	return Read(strings.NewReader(source))
}

// ReadFile reads the SML source of a configuration file,
// parses it, and returns the etc instance.
func ReadFile(filename string) (*Etc, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, failure.Annotate(err, "cannot read file '%s'", filename)
	}
	return ReadString(string(source))
}

// HasPath checks if the configurations has the defined path
// regardles of the value or possible subconfigurations.
func (e *Etc) HasPath(path string) bool {
	fullPath := makeFullPath(path)
	changer := e.values.At(fullPath...)
	return changer.Error() == nil
}

// Do iterates over the children of the given path and executes
// the function f with that path.
func (e *Etc) Do(path string, f func(p string) error) error {
	fullPath := makeFullPath(path)
	changer := e.values.At(fullPath...)
	if changer.Error() != nil {
		return changer.Error()
	}
	kvs, err := changer.List()
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		p := pathToString(append(fullPath, kv.Key))
		err := f(p)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValueAsString retrieves the string value at a given path. If it
// doesn't exist the default value dv is returned.
func (e *Etc) ValueAsString(path, dv string) string {
	value := e.valueAt(path)
	return defaulter.AsString(value, dv)
}

// ValueAsBool retrieves the bool value at a given path. If it
// doesn't exist the default value dv is returned.
func (e *Etc) ValueAsBool(path string, dv bool) bool {
	value := e.valueAt(path)
	return defaulter.AsBool(value, dv)
}

// ValueAsInt retrieves the int value at a given path. If it
// doesn't exist the default value dv is returned.
func (e *Etc) ValueAsInt(path string, dv int) int {
	value := e.valueAt(path)
	return defaulter.AsInt(value, dv)
}

// ValueAsFloat64 retrieves the float64 value at a given path. If it
// doesn't exist the default value dv is returned.
func (e *Etc) ValueAsFloat64(path string, dv float64) float64 {
	value := e.valueAt(path)
	return defaulter.AsFloat64(value, dv)
}

// ValueAsTime retrieves the string value at a given path and
// interprets it as time with the passed format. If it
// doesn't exist the default value dv is returned.
func (e *Etc) ValueAsTime(path, format string, dv time.Time) time.Time {
	value := e.valueAt(path)
	return defaulter.AsTime(value, format, dv)
}

// ValueAsDuration retrieves the duration value at a given path.
// If it doesn't exist the default value dv is returned.
func (e *Etc) ValueAsDuration(path string, dv time.Duration) time.Duration {
	value := e.valueAt(path)
	return defaulter.AsDuration(value, dv)
}

// Split produces a subconfiguration below the passed path.
// The last path part will be the new root, all values below
// that configuration node will be below the created root.
// In case of an invalid path an empty configuration will
// be returned as default.
func (e *Etc) Split(path string) (*Etc, error) {
	if !e.HasPath(path) {
		// Path not found, return empty configuration.
		return ReadString("{etc}")
	}
	fullPath := makeFullPath(path)
	values, err := e.values.CopyAt(fullPath...)
	if err != nil {
		return nil, failure.Annotate(err, "cannot split configuration")
	}
	values.At(fullPath[len(fullPath)-1:]...).SetKey("etc")
	es := &Etc{
		values: values,
	}
	return es, nil
}

// Dump creates a map of paths and their values to apply
// them into other configurations.
func (e *Etc) Dump() (Application, error) {
	appl := Application{}
	err := e.values.DoAllDeep(func(ks []string, v string) error {
		if len(ks) == 1 {
			// Continue on root element.
			return nil
		}
		path := strings.Join(ks[1:], "/")
		appl[path] = v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return appl, nil
}

// Apply creates a new configuration by adding of overwriting
// the passed values. The keys of the map have to be slash
// separated configuration paths without the leading "etc".
func (e *Etc) Apply(appl Application) (*Etc, error) {
	ec := &Etc{
		values: e.values.Copy(),
	}
	for path, value := range appl {
		fullPath := makeFullPath(path)
		_, err := ec.values.Create(fullPath...).SetValue(value)
		if err != nil {
			return nil, failure.Annotate(err, "cannot apply changes")
		}
	}
	return ec, nil
}

// Write writes the configuration as SML to the passed target.
// If prettyPrint is true the written SML is indented and has
// linebreaks.
func (e *Etc) Write(target io.Writer, prettyPrint bool) error {
	// Build the nodes tree.
	builder := sml.NewNodeBuilder()
	depth := 0
	err := e.values.DoAllDeep(func(ks []string, v string) error {
		doDepth := len(ks)
		tag := ks[doDepth-1]
		for i := depth; i > doDepth; i-- {
			builder.EndTagNode()
		}
		switch {
		case doDepth > depth:
			builder.BeginTagNode(tag)
			builder.TextNode(v)
			depth = doDepth
		case doDepth == depth:
			builder.EndTagNode()
			builder.BeginTagNode(tag)
			builder.TextNode(v)
		case doDepth < depth:
			builder.EndTagNode()
			builder.BeginTagNode(tag)
			builder.TextNode(v)
			depth = doDepth
		}
		return nil
	})
	if err != nil {
		return err
	}
	for i := depth; i > 0; i-- {
		builder.EndTagNode()
	}
	root, err := builder.Root()
	if err != nil {
		return err
	}
	// Now write the node structure.
	wp := sml.NewStandardSMLWriter()
	wctx := sml.NewWriterContext(wp, target, prettyPrint, "   ")
	return sml.WriteSML(root, wctx)
}

// String implements the fmt.Stringer interface.
func (e *Etc) String() string {
	return fmt.Sprintf("%v", e.values)
}

// valueAt retrieves and encapsulates the value
// at a given path.
func (e *Etc) valueAt(path string) *value {
	fullPath := makeFullPath(path)
	changer := e.values.At(fullPath...)
	return &value{fullPath, changer}
}

// postProcess replaces templates formated [path||default]
// with values found at that path or the default.
func (e *Etc) postProcess() error {
	re := regexp.MustCompile(`\[.+(||.+)\]`)
	// Find all entries with template.
	changers := e.values.FindAll(func(k, v string) (bool, error) {
		return re.MatchString(v), nil
	})
	// Change the template.
	for _, changer := range changers {
		value, err := changer.Value()
		if err != nil {
			return err
		}
		found := re.FindString(value)
		// Look for default value.
		sourceDefault := strings.SplitN(found[1:len(found)-1], "||", 2)
		defaultValue := found
		if len(sourceDefault) > 1 {
			defaultValue = sourceDefault[1]
		}
		// Check if source is environment variable or path.
		substitute := ""
		if strings.HasPrefix(sourceDefault[0], "$") {
			if envValue, ok := os.LookupEnv(sourceDefault[0][1:]); ok {
				substitute = envValue
			} else {
				substitute = defaultValue
			}
		} else {
			substitute = e.ValueAsString(sourceDefault[0], defaultValue)
		}
		replaced := strings.Replace(value, found, substitute, -1)
		_, err = changer.SetValue(replaced)
		if err != nil {
			return err
		}
	}
	return nil
}

//--------------------
// CONTEXT
//--------------------

// NewContext returns a new context that carries a configuration.
func NewContext(ctx context.Context, cfg *Etc) context.Context {
	return context.WithValue(ctx, etcKey, cfg)
}

// FromContext returns the configuration stored in ctx, if any.
func FromContext(ctx context.Context) (*Etc, bool) {
	cfg, ok := ctx.Value(etcKey).(*Etc)
	return cfg, ok
}

//--------------------
// HELPERS
//--------------------

// makeFullPath creates the full path out of a string.
func makeFullPath(path string) []string {
	parts := stringex.SplitMap(path, "/", func(p string) (string, bool) {
		if p == "" {
			return "", false
		}
		return strings.ToLower(p), true
	})
	return append(etcRoot, parts...)
}

// fullPathToString returns the path in a filesystem like notation.
func fullPathToString(path []string) string {
	return "/" + strings.Join(path, "/")
}

// pathToString returns the path in a filesystem like notation but
// with the leading slash and 'etc'.
func pathToString(path []string) string {
	return strings.Join(path[1:], "/")
}

// EOF
