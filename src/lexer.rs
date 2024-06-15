use crate::types::{Lexer, Position, Token, TokenTypes, TokenValue};
use std::num::ParseIntError;

fn new_lexer<'a>(filename: String, text: String) -> Lexer {
    let bytes = String::from(text).as_bytes();
    let mut lexer = Lexer {
        fn_name: filename.clone(),
        text: bytes,
        pos: new_position(0, 0, -1, filename, bytes),
        current_char: 0,
    };

    lexer.advance();
    return lexer;
}

fn new_position(idx: usize, ln: i32, col: i32, fn_name: String, ftxt: &[u8]) -> Position {
    return Position {
        idx,
        ln,
        col,
        fn_name,
        ftxt,
    };
}

impl Lexer {
    fn advance(&mut self) {
        self.pos.advance(self.current_char);
        if self.pos.idx < self.text.len() {
            self.current_char = self.text[self.pos.idx];
        } else {
            self.current_char = ' ';
        }
    }
}

impl Lexer {
    fn make_tokens(&mut self) -> Vec<&Token> {
        let mut tokens: Vec<&Token> = vec![];

        while self.current_char != None {
            match self.current_char {
                'a'..='z' | 'A'..='Z' => tokens.push(self.make_identifier()),
                '0'..='9' => tokens.push(self.make_number()),
                ' ' | '\t' | '\r' => {}
                '\n' | ';' => tokens.push(new_token(
                    TokenTypes::Newline,
                    None,
                    self.pos.clone(),
                    self.pos.clone(),
                )),
                '"' => tokens.push(self.make_letter()),
                '+' => {}
                '-' => {}
                '{' => {}
                '}' => {}
                '/' => {}
                '(' => {}
                ')' => {}
                '[' => {}
                ']' => {}
                '^' => {}
                '!' => {}
                '=' => {}
                '<' => {}
                '>' => {}
                ',' => {}
                '.' => {}
                '&' => {}
                '*' => {}
                _ => {
                    let pos_start = self.pos.clone();
                    let char = self.current_char;
                    self.advance();
                    return;
                }
            }
        }

        return tokens;
    }
}

fn new_token<'a>(
    token_type: TokenTypes,
    value: Option<TokenValue>,
    pos_start: Position,
    pos_end: Position,
) -> &'a Token {
    return &Token {
        r#type: token_type,
        value,
        pos_start,
        pos_end,
    };
}

impl Lexer {
    fn make_number(&mut self) -> &Token {
        let mut num_string = String::from("");
        let mut dot_count = 0;
        let pos_start = self.pos.clone();

        while let Some(ch) = self.current_char {
            if is_digit(ch) || ch == '.' {
                if ch == '.' {
                    if dot_count == 1 {
                        break;
                    }
                    dot_count += 1;
                    num_string.push('.')
                } else {
                    num_string.push(ch)
                }
                self.advance()
            } else {
                break;
            }
        }

        let mut pos_end = self.pos.clone();
        pos_end.col -= 1;
        pos_end.idx -= 1;

        if dot_count == 0 {
            let parsed_float: Result<f64, _> = num_string.parse();
            return new_token(
                TokenTypes::Float,
                Some(TokenValue::Float(parsed_float.unwrap())),
                pos_start,
                pos_end,
            );
        }

        let parsed_int: Result<i64, _> = num_string.parse();
        new_token(
            TokenTypes::Int,
            Some(TokenValue::Integer(parsed_int.unwrap())),
            pos_start,
            pos_end,
        )
    }
}
impl Lexer {
    fn make_identifier(&mut self) -> &Token {

    }
}

impl Lexer {
    fn make_string(&mut self) -> &Token {

    }
}

fn is_digit(ch: char) -> bool {
    ch.is_digit(10)
}
