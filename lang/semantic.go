package lang

// Symbol represents the program entities that we would want to track via the
// symbol table
type Symbol interface {
	Name() string
	String() string
	setScope(Scope)
}

// baseSymbol is the base implementation for a Symbol, to be embedded
type baseSymbol struct {
	name  string // name of the symbol
	scope Scope  // all symbols track their scope
}

// Name returns the name of the symbol
func (s baseSymbol) Name() string { return s.name }

// setScope sets the scope field of a symbol
func (s baseSymbol) setScope(scope Scope) { s.scope = scope }

func (s baseSymbol) String() string { return s.Name() }

// VarSymbol is symbol that represents a variable (using an identifier)
type VarSymbol struct{ baseSymbol }

// Scope refers to the scope that is used to track symbols from the program
type Scope interface {
	ScopeName() string
	EnclosingScope() (Scope, bool)      // gets the parent scope if available
	Resolve(name string) (Symbol, bool) // lookup scopenames
	// private
	define(Symbol) // define symbols in this scope
}

// DefineSymbol defines a symbol in the given scope, adding it into the scope
// as well as setting the symbol's scope to this scope.
func DefineSymbol(symbol Symbol, scope Scope) {
	scope.define(symbol)
	symbol.setScope(scope)
}

// baseScope implements most of the base implementation of Scopes in went
// NOTE: baseScope is not a complete implementation of Scope (Does not implement
// ScopeName), should be embedded
type baseScope struct {
	enclosingScope Scope
	symbols        map[string]Symbol
}

func (s *baseScope) EnclosingScope() (Scope, bool) {
	if s.enclosingScope != nil {
		return nil, false
	}
	return s.enclosingScope, true
}

func (s *baseScope) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.symbols[name]
	if ok {
		return symbol, ok
	}
	es, ok := s.EnclosingScope()
	if ok {
		return es.Resolve(name)
	}
	return nil, false
}

// puts the symbol in the symbols map, not meant to be directly called
func (s *baseScope) define(symbol Symbol) { s.symbols[symbol.Name()] = symbol }

// GlobalScope is the top level scope in the program, it has no enclosing scope
type GlobalScope struct{ baseScope }

// ScopeName returns global
func (s *GlobalScope) ScopeName() string { return "global" }

// NewGlobalScope returns a new global scope
func NewGlobalScope() *GlobalScope { return &GlobalScope{baseScope{}} }

// LocalScope is any local scope that is created by the program via blocks
// These blocks are enclosed in '{' '}'
type LocalScope struct{ baseScope }

// ScopeName returns local
func (s *LocalScope) ScopeName() string { return "local" }

// NewLocalScope returns a new local scope, that has its parent set
func NewLocalScope(parent Scope) *LocalScope {
	return &LocalScope{baseScope{enclosingScope: parent}}
}

// FunctionSymbol is a Symbol that also has a Scope, does not use baseScope's
// implementation
type FunctionSymbol struct {
	baseSymbol
	enclosingScope Scope
	formalArgs     map[string]Symbol
	funcBlock      Node // holds the node where the function block is
}

// TODO: Override baseSymbol's Name()
