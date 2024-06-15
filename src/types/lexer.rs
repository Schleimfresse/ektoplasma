use crate::types::{Position, TokenTypes};

pub struct Lexer<'a> {
    pub fn_name: String,
    pub text: &'a [u8],
    pub pos: Position<'a>,
    pub current_char: char,
}

pub struct Token<'a> {
    pub r#type: TokenTypes,
    pub value:  Option<TokenValue>,
    pub pos_start: Position<'a>,
    pub pos_end: Position<'a>,
}

pub enum TokenValue {
    Float(f64),
    Integer(i64),
    Str(String),
}