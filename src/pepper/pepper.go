// Package pepper outputs rpc service descriptions in Go code.
// It runs as a plugin for the Go protocol buffer compiler plugin.
// It is linked in to protoc-gen-go.
package pepper

import (
	"fmt"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func init() {
	generator.RegisterPlugin(new(pepperGen))
}

// pepperGen is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for pepper support.
type pepperGen struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "pepper".
func (g *pepperGen) Name() string {
	return "pepper"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	pepperPkg string
)

// Init initializes the plugin.
func (g *pepperGen) Init(gen *generator.Generator) {
	g.gen = gen
	pepperPkg = generator.RegisterUniquePackageName("pepper", nil)
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *pepperGen) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *pepperGen) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *pepperGen) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *pepperGen) Generate(file *generator.FileDescriptor) {

	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

    g.gen.AddImport("bytes")
    g.gen.AddImport("compress/gzip")
    g.gen.AddImport("io/ioutil")
    g.gen.AddImport("github.com/golang/protobuf/protoc-gen-go/descriptor")

	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}

    g.generateVar(file)
    g.generateInit(file)
}

// GenerateImports generates the import declaration for this file.
func (g *pepperGen) GenerateImports(file *generator.FileDescriptor) {
}

func (g *pepperGen) generateVar(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

    g.P()
    g.P("var (")
	for _, service := range file.FileDescriptorProto.Service {
        origServName := service.GetName()
        servName := generator.CamelCase(origServName)
        g.P("    ", servName, "_Desc *descriptor.ServiceDescriptorProto")
	}
    g.P(")")
    g.P()
    g.P("var (")
    g.P("    ", file.VarName(), "_service_descs = make(map[string]*descriptor.ServiceDescriptorProto)")
    g.P(")")
    g.P()
}

func (g *pepperGen) generateInit(file *generator.FileDescriptor) {
    g.P()
    g.P("func init() {")
    g.P("    r, err := gzip.NewReader(bytes.NewReader(", file.VarName(), "))")
    g.P("    if err != nil {")
    g.P("        panic(err.Error())")
    g.P("    }")
    g.P("    defer r.Close()")
    g.P()
    g.P("    b, err := ioutil.ReadAll(r)")
    g.P("    if err != nil {")
    g.P("        panic(err.Error())")
    g.P("    }")
    g.P()
    g.P("    file_desc := new(descriptor.FileDescriptorProto)")
    g.P("    if err := proto.Unmarshal(b, file_desc); err != nil {")
    g.P("        panic(err.Error())")
    g.P("    }")
    g.P()
    g.P("    for _, service := range file_desc.Service {")
    g.P("        ", file.VarName(), "_service_descs[service.GetName()] = service")
    g.P("    }")
    for _, service := range file.FileDescriptorProto.Service {
        origServName := service.GetName()
        servName := generator.CamelCase(origServName)
        g.P("    ", servName, "_Desc = ", file.VarName(), "_service_descs[\"", servName, "\"]")
    }
    g.P("}")
    g.P()
}

// generateService generates all the code for the named service.
func (g *pepperGen) generateService(file *generator.FileDescriptor, service *desc.ServiceDescriptorProto, index int) {
    path := fmt.Sprintf("6,%d", index) // 6 means service.
    _ = path

    origServName := service.GetName()
	servName := generator.CamelCase(origServName)

	g.P()

	// Server interface.
	g.P("type ", servName, " interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		g.P(g.genMethod(servName, method))
	}
	g.P("}")
	g.P()
}

func (g *pepperGen) genMethod(servName string, method *desc.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())
	return methName + "(uint32, *" + inType + ", *" + outType + ") error"
}
