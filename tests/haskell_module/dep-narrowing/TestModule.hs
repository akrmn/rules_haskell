module TestModule(bar, foo2) where

import TestLibModule (foo)
import TestLibModule2 (foo2)

bar :: Int
bar = 2 * foo + foo2
