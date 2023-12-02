# Doubly Linked List move function

## move src before dst


| src\dst                    | not in list | the only one | the first but not the last           | the last but not the first                       | in middle             |
| -------------------------- | ----------- | ------------ | ------------------------------------ | ------------------------------------------------ | --------------------- |
| not in list                | x           | x            | x                                    | x                                                | x                     |
| the only one               | x           | x            | x                                    | x                                                | x                     |
| the first but not the last | x           | x            | x                                    | root.header(src.next)<br /><br /> root.tail(src) | root.header(src.next) |
| the last but not the first | x           | x            | root.header(src) root.tail(src.prev) | x                                                | root.tail(src.prev)   |
| in middle                  | x           | x            | root.header(src)                     | x                                                | x                     |

## move src after dst


| src\dst                    | not in list | the only one | the first but not the last           | the last but not the first                       | in middle             |
| -------------------------- | ----------- | ------------ | ------------------------------------ | ------------------------------------------------ | --------------------- |
| not in list                | x           | x            | x                                    | x                                                | x                     |
| the only one               | x           | x            | x                                    | x                                                | x                     |
| the first but not the last | x           | x            | x                                    | root.header(src.next)<br /><br /> root.tail(src) | root.header(src.next) |
| the last but not the first | x           | x            | root.header(src) root.tail(src.prev) | x                                                | root.tail(src.prev)   |
| in middle                  | x           | x            | x                                    | root.tail(src)                                   | x                     |
