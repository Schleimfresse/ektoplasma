mod lexer;
mod nodes;
mod position;
pub use lexer::Lexer;
pub use position::Position;
pub use lexer::Token;
pub use lexer::TokenValue;
pub type Binary = i8;

pub enum TokenTypes {
    Int,
    Float,
    String,
    Identifier,
    Keyword,
    Plus,
    Minus,
    Div,
    Eq,
    Lparen,
    Rparen,
    Lsquare,
    Rsquare,
    Lbrace,
    Rbrace,
    Pow,
    EE,
    NE,
    LT,
    GT,
    LTE,
    GTE,
    EOF,
    Comma,
    Newline,
    Arrow,
    Dot,
    And,
    Star,
}