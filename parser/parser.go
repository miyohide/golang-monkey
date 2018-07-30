package parser

import (
	"fmt"

	"github.com/miyohide/monkey/ast"
	"github.com/miyohide/monkey/lexer"
	"github.com/miyohide/monkey/token"
)

// Parser は字句解析器インスタンスへのポインタと2つのトークンを持つ型
type Parser struct {
	// 字句解析器インスタンスへのポインタ
	l *lexer.Lexer

	// 現在のトークン
	curToken token.Token
	// 次のトークン
	peekToken token.Token
	// 入力に想定外のトークンに遭遇したときに入れる
	errors []string
}

// New は新しく字句解析器を作成し、2つのトークンを読み込んだParserを返す
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 2つトークンを読み込んでcurTokenとpeekTokenの両方をセット
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken はcurTokenとpeekTokenをセットするヘルパーメソッド
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram はプログラムをParseする関数
func (p *Parser) ParseProgram() *ast.Program {
	// ASTのルートノードを生成する
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// TODO セミコロンに遭遇するまで式を読み飛ばしてしまっている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

// Errors は現在のエラーの配列を返す
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
