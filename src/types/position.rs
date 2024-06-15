pub struct Position<'a> {
    pub idx: usize,
    pub ln: i32,
    pub col: i32,
    pub fn_name: String,
    pub ftxt: &'a [u8],
}

impl Clone for Position {
    fn clone(&self) -> Self {
        Self {
            idx: self.idx,
            ln: self.ln,
            col: self.col,
            fn_name: self.fn_name.clone(),
            ftxt: self.ftxt.clone(),
        }
    }
}

impl Position {
    pub fn advance(&mut self, current_char: u8) -> &mut Position {
        self.idx += 1;
        self.col += 1;

        if current_char == u8::try_from('\n').unwrap() {
            self.ln += 1;
            self.col = 0
        }
        return self;
    }
}
