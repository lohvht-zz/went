package lang

// Symbol represents the program entities that we would want to track via the
// symbol table
type Symbol interface {
	Name() string
	SetScope(Scope)
	String() string
}

// baseSymbol is the base implementation for a Symbol, to be embedded
type baseSymbol struct {
	name  string // name of the symbol
	scope Scope  // all symbols track their scope
}

// Name returns the name of the symbol
func (s baseSymbol) Name() string { return s.name }

// SetScope sets the scope field of a symbol
func (s baseSymbol) SetScope(scope Scope) { s.scope = scope }

func (s baseSymbol) String() string { return s.Name() }

// TypeSymbol is a symbol that represents any given type, self-defined or built-in
type TypeSymbol struct{ baseSymbol }

// Built-in Types Symbols
var (
	intType   = TypeSymbol{baseSymbol{name: "int"}}
	floatType = TypeSymbol{baseSymbol{name: "float"}}
	listType  = TypeSymbol{baseSymbol{name: "list"}}
	mapType   = TypeSymbol{baseSymbol{name: "map"}}
	nullType  = TypeSymbol{baseSymbol{name: "null"}}
	boolType  = TypeSymbol{baseSymbol{name: "bool"}}

	DefaultTypeMap = map[string]TypeSymbol{
		"int":   intType,
		"float": floatType,
		"list":  listType,
		"map":   mapType,
		"null":  nullType,
		"bool":  boolType,
	}
)

// Scope refers to the scope that is used to track symbols from the program
type Scope interface {
	ParentScope() Scope                 // gets the parent scope
	Define(Symbol)                      // define symbols in this scope
	Resolve(name string) (Symbol, bool) // lookup scopenames
}

// baseScope implements most of the important Scope logic, to be embedded within
// other implementations of Scope
type baseScope struct {
	parentScope Scope
	symbols     map[string]Symbol
}

func newBaseScope(parent Scope) baseScope {
	return baseScope{parentScope: parent, symbols: map[string]Symbol{}}
}

func (bs baseScope) ParentScope() Scope {
	return bs.parentScope
}

func (bs baseScope) Define(symb Symbol) {
	bs.symbols[symb.Name()] = symb
	symb.SetScope(bs)
}

func (bs baseScope) Resolve(name string) (Symbol, bool) {
	s, ok := bs.symbols[name]
	if ok {
		return s, true
	}
	return nil, ok
}

// GlobalScope is the global context that should be accessible by other scopes
type GlobalScope struct{ baseScope }

// NewGlobalScope returns a pointer to a GlobalScope
func NewGlobalScope() *GlobalScope {
	return &GlobalScope{baseScope: newBaseScope(nil)}
}

// LocalScope is an accessible inner scope
type LocalScope struct{ baseScope }

// NewLocalScope returns a pointer to a LocalScope
func NewLocalScope(parent Scope) *LocalScope {
	return &LocalScope{baseScope: newBaseScope(parent)}
}

// SymbolTable holds the symbol for 1 run of the interpreter
type SymbolTable struct{ globals *GlobalScope }

// NewSymbolTable initialises the built-in types
func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{globals: NewGlobalScope()}
	st.initTypeSystem()
	return st
}

// initTypeSystem initialises the built-in types that went supports
func (symbtab *SymbolTable) initTypeSystem() {
	for _, v := range DefaultTypeMap {
		symbtab.globals.Define(v)
	}
}
