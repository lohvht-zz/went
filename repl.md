# REPL Specs: (Based on python interactive shell

1) Left curly bracket => able to put multiple lines, enter keys pressed after right curly brackets inserted will evaluate
2) Every enter inserts a newline character '\n', from single line expressions to curly brackets
3) If a line should be interpreted as a single line. but you want to break it into multiple lines, put a '\' at the end of the line
4) 1st line/Single line inputs will have a prefix "[went]> " while multiple line inputs will have subsequent lines be "....... "

    ==> Multiline input condition
    1) Inside a "block", a "list", or a "function" (see that we're inside a "(", "{", "[")
        1.1) Easy way to exit out of this is to just put right bracket of any type (Does not have to match)
        1.2) Multiple lines in this condition will insert a '\n' character to the entered line
    2) If the last character of the line is a "\"
        2.2) Multiple lines in this condition will insert a space(" ") instead

Every other newline will have a '\n' newline entered
