# ============================================================
# PYTHON LEXER TORTURE FILE
# PART 1/4
#
# Target:
#   Python 3.14+
#
# Includes:
#   - imports
#   - literals
#   - strings
#   - f-strings
#   - t-strings
#   - raw strings
#   - bytes
#   - unicode
#   - operators
# ============================================================


from __future__ import annotations


import asyncio
import cmath
import datetime
import functools
import itertools
import json
import math
import pathlib
import re
import statistics
import sys
from collections import Counter, defaultdict
from dataclasses import dataclass
from typing import (
    Any,
    Callable,
    ClassVar,
    Generic,
    Iterable,
    Iterator,
    Literal,
    Protocol,
    Self,
    TypeVar,
)


# ============================================================
# Comments and whitespace
# ============================================================


# Single line comment

x = 1  # trailing comment


# ============================================================
# Numeric literals
# ============================================================


INTEGER_DECIMAL = 123456789

INTEGER_UNDERSCORE = 123_456_789

INTEGER_BINARY = 0b101010101

INTEGER_OCTAL = 0o755

INTEGER_HEX = 0xDEADBEEF


FLOAT_SIMPLE = 123.456

FLOAT_EXPONENT = 1.23e45

FLOAT_NEGATIVE_EXPONENT = 9.99e-10

FLOAT_UNDERSCORE = 1_000.000_001


COMPLEX_SIMPLE = 3j

COMPLEX_FULL = 10 + 20j


# ============================================================
# Operators and delimiters
# ============================================================


OPERATORS = [
    "+",
    "-",
    "*",
    "**",
    "/",
    "//",
    "%",
    "@",

    "<<",
    ">>",

    "&",
    "|",
    "^",
    "~",

    ":=",

    "<",
    ">",
    "<=",
    ">=",

    "==",
    "!=",

    "(",
    ")",
    "[",
    "]",
    "{",
    "}",

    ",",
    ":",
    ".",
    ";",

    "->",
]


operator_expression = (
    (
        +1
        - 2
        * 3
        / 4
        // 5
        % 6
    )
    ** 2
)


bit_expression = (
    (
        0b1010
        &
        0b1111
    )
    |
    (
        0b0101
    )
)


shift_expression = (
    1
    << 8
    >>
    2
)


matrix_operator = (
    left
    @
    right
)


walrus_expression = (
    value := calculate()
)


# ============================================================
# Strings
# ============================================================


EMPTY = ""

SINGLE = 'single'

DOUBLE = "double"


ESCAPED = (
    "\n"
    "\t"
    "\r"
    "\\"
    "\""
    "\'"
    "\x41"
    "\u1234"
    "\U0001F600"
)


RAW = r"""
C:\Users\someone\folder\file.txt
\n is not escaped
"""


TRIPLE_SINGLE = '''
hello
world
'''


TRIPLE_DOUBLE = """
hello
world
"""


BYTES = b"hello"

RAW_BYTES = rb"C:\windows\path"


BYTES_ESCAPES = (
    b"\x00"
    b"\xff"
)


# Adjacent strings

JOINED = (
    "hello "
    "world"
)


# ============================================================
# Unicode identifiers
# ============================================================


π = 3.14159

café = "coffee"

naïve = True


unicode_identifier_test = (
    π,
    café,
    naïve,
)


# ============================================================
# f-string torture
# ============================================================


name = "Lexer"
number = 42
value = 100


F_BASIC = f"hello {name}"


F_CONVERSION = (
    f"{name!r}"
    f"{name!s}"
    f"{name!a}"
)


F_FORMAT = (
    f"{number:08x}"
)


F_DEBUG = (
    f"{name=}"
    f"{number=}"
)


F_ARITHMETIC = f"""
sum:
    {1 + 2}

product:
    {3 * 4}

power:
    {2 ** 16}
"""


F_ATTRIBUTE = (
    f"{object.value}"
)


F_CALL = (
    f"{function(arg1, arg2)}"
)


F_INDEX = (
    f"{items[0]}"
)


F_SLICE = (
    f"{items[1:10:2]}"
)


F_DICT_ACCESS = (
    f"{config['database']['host']}"
)


# ============================================================
# f-string nested expressions
# ============================================================


F_NESTED_EXPRESSION = f"""
Balance:
{
    {
        "uk": "£"
    }["uk"]
    +
    account.balance
    :
    ,.2f
}
"""


F_NESTED_FORMAT = (
    f"{math.pi:{width}.{precision}f}"
)


F_ESCAPED_BRACES = (
    f"""
    {{
        literal
    }}

    {{{value}}}
    """
)


F_NESTED_FSTRING = (
    f"{f'{name}:{number:04d}'}"
)


# ============================================================
# Raw f strings
# ============================================================


RAW_F = rf"""
Path:
    C:\Users\{name}\Documents

Regex:
    \d+\w+\s+
"""


RAW_TRIPLE_F = rf"""
before

{
    value
}

after
"""


# ============================================================
# t-string (Python 3.14)
# ============================================================


T_SIMPLE = t"hello {name}"


T_COMPLEX = t"""
User:
    {name}

Balance:
    {account.balance:,.2f}
"""


T_WITH_ESCAPED_BRACES = t"""
{{ literal }}

{{{value}}}
"""


# ============================================================
# Multiline continuation
# ============================================================


continued = 1 + \
    2 + \
    3


continued_string = "hello " \
                   "world"


# ============================================================
# End Part 1
# ============================================================

# ============================================================
# PYTHON LEXER TORTURE FILE
# PART 2/4
#
# Advanced syntax:
#   - decorators
#   - signatures
#   - generics
#   - classes
#   - match
#   - async
#   - exceptions
# ============================================================


# ============================================================
# Type torture
# ============================================================


T = TypeVar("T")
K = TypeVar("K")
V = TypeVar("V")


SimpleType = int | str | None


ComplexType = dict[
    str,
    list[
        tuple[
            int,
            str | None,
        ]
    ],
]


NestedCallable = Callable[
    [
        int,
        str,
        list[float],
    ],
    dict[str, Any],
]


# ============================================================
# Protocols
# ============================================================


class SupportsMagic(Protocol):

    def magic(
        self,
        value: int,
    ) -> str:
        ...


# ============================================================
# Decorator torture
# ============================================================


def decorator(
    value,
    *,
    enabled=True,
):

    def wrapper(function):

        function.metadata = value

        return function

    return wrapper



@decorator(
    f"{name}:{number}",
    enabled=(
        lambda x: x > 10
    )(100),
)
def decorated_function():

    return True



@(
    decorator_factory(
        "nested",
        enabled=True,
    )
)
def another_decorated_function():

    pass



@staticmethod
@classmethod
@decorator("stacked")
def stacked():

    pass



# ============================================================
# Function signatures from hell
# ============================================================


def signature_torture(
    positional_only,
    /
    ,
    normal,
    default=value,
    *args,
    keyword_only,
    optional=123,
    **kwargs,
) -> dict[str, Any]:

    return {
        "args": args,
        "kwargs": kwargs,
    }



def annotations_everywhere(
    value: list[
        tuple[
            int,
            str,
        ]
    ],
    callback: Callable[
        [int],
        str,
    ] = lambda x: str(x),
) -> tuple[
    bool,
    str,
]:

    return True, callback(10)



# ============================================================
# PEP 695 generic syntax (3.12+)
# ============================================================


class Cache[K, V]:

    def __init__(
        self,
        key: K,
        value: V,
    ):

        self.key = key
        self.value = value



def identity[T](
    value: T,
) -> T:

    return value



class TreeNode[T]:

    def __init__(
        self,
        value: T,
        children: list[TreeNode[T]],
    ):

        self.value = value
        self.children = children



# ============================================================
# Nested classes
# ============================================================


class Outer:

    class Middle:

        class Inner:

            def value(
                self,
            ):

                return (
                    {
                        "deep": [
                            {
                                "value": True
                            }
                        ]
                    }
                )


    def method(self):

        class Local:

            def run(self):

                return "local"

        return Local()



# ============================================================
# Properties and descriptors
# ============================================================


class Descriptor:

    def __get__(
        self,
        instance,
        owner,
    ):

        return (
            instance,
            owner,
        )



class PropertyExample:

    attribute = Descriptor()


    @property
    def value(
        self,
    ) -> int:

        return 42


    @value.setter
    def value(
        self,
        new_value: int,
    ):

        self._value = new_value



# ============================================================
# Pattern matching torture
# ============================================================


def pattern_matching(
    token,
):

    match token:


        case None:

            return "none"


        case 0 | 1 | 2:

            return "small"


        case int(value):

            return value


        case str(value) if len(value) > 10:

            return "long string"


        case [
            first,
            second,
            *remaining,
        ]:

            return (
                first,
                second,
                remaining,
            )


        case {
            "name": username,
            "age": age,
            **extra,
        }:

            return (
                username,
                age,
                extra,
            )


        case Token(
            kind="NUMBER",
            value=value,
        ) if value > 100:

            return "large number"


        case _:

            return "unknown"



# ============================================================
# Match with complicated guards
# ============================================================


def guarded_match(value):

    match value:

        case x if (
            x is not None
            and x > 10
            and (
                lambda y: y
            )(True)
        ):

            return x


        case _:

            return None



# ============================================================
# Exception torture
# ============================================================


def exception_torture():

    try:

        raise ValueError(
            "broken"
        )


    except (
        ValueError,
        TypeError,
    ) as error:

        print(error)


    except* RuntimeError:

        pass


    finally:

        cleanup()



# ============================================================
# Async torture
# ============================================================


async def async_torture(
    source: AsyncIterable[int],
):

    values = []


    async for item in source:

        values.append(item)


    return values



async def async_generator_torture():

    for number in range(10):

        await asyncio.sleep(
            0
        )

        yield number



async def async_with_torture():

    async with (
        AsyncContext()
        as first,
        AsyncContext()
        as second,
    ):

        return (
            first,
            second,
        )



# ============================================================
# Context manager torture
# ============================================================


def with_torture():

    with (
        open("a.txt") as first,
        open("b.txt") as second,
        open("c.txt") as third,
    ):

        return (
            first,
            second,
            third,
        )



# ============================================================
# yield torture
# ============================================================


def generator_torture():

    yield from (
        x
        for x in range(100)
        if x % 2 == 0
    )


def nested_generator():

    yield (
        (
            yield_value
        )
        for yield_value in range(5)
    )


# ============================================================
# End Part 2
# ============================================================

# ============================================================
# PYTHON LEXER TORTURE FILE
# PART 3/4
#
# Maximum expression nesting
# ============================================================


# ============================================================
# f-string expression abuse
# ============================================================


f_comment_expression = f"""
{
    (
        value
        # comment inside expression
        +
        number
    )
}
"""


f_string_inside_expression = f"""
{
    (
        "hello"
        +
        f"{name}:{number}"
    )
}
"""


f_triple_string_inside_expression = f"""
{
    '''
    triple
    quoted
    string
    '''
}
"""


f_dict_inside_expression = f"""
{
    {
        "key":
            {
                "nested": True
            }
    }["key"]["nested"]
}
"""


f_lambda_inside_expression = f"""
{
    (
        lambda x:
            x * x
    )(10)
}
"""


f_comprehension_inside_expression = f"""
{
    [
        x * x
        for x in range(20)
        if x % 2
    ]
}
"""


f_nested_everything = f"""
Result:

{
    (
        lambda data:
            {
                key:
                    value
                for key, value in data.items()
                if (
                    value
                    and
                    value > 10
                )
            }
    )(
        {
            "a": 1,
            "b": 20,
        }
    )
}
"""


# ============================================================
# Nested comprehensions
# ============================================================


mega_list = [
    [
        [
            (
                a,
                b,
                c,
            )
            for c in range(5)
        ]
        for b in range(5)
    ]
    for a in range(5)
]



mega_dict = {
    outer_key:
        {
            inner_key:
                [
                    value
                    for value in values
                    if (
                        value
                        >
                        threshold
                    )
                ]
            for inner_key, values in inner.items()
        }
    for outer_key, inner in data.items()
}



mega_generator = (
    (
        x,
        y,
        z,
    )
    for x in range(10)
    for y in range(10)
    for z in range(10)
    if (
        x
        +
        y
        +
        z
    )
)



# ============================================================
# Lambda nightmares
# ============================================================


lambda_chain = (
    lambda x:
        (
            lambda y:
                (
                    lambda z:
                        x + y + z
                )(3)
        )(2)
)(1)



lambda_with_defaults = (
    lambda
        a=1,
        b=2,
        *args,
        **kwargs:
            (
                a,
                b,
                args,
                kwargs,
            )
)



# ============================================================
# Assignment torture
# ============================================================


simple_assignment = 123


multiple_assignment = (
    a,
    b,
    c,
) = (
    1,
    2,
    3,
)


star_assignment = (
    first,
    *middle,
    last,
) = range(100)



attribute_assignment = (
    object.attribute
    ==
    value
)



subscript_assignment = (
    matrix[0][1]
    ==
    42
)



annotated_assignment: list[int] = [
    1,
    2,
    3,
]


# ============================================================
# Global / nonlocal
# ============================================================


global_value = 0


def outer_function():

    outer_value = 10


    def inner_function():

        nonlocal outer_value

        global global_value

        outer_value += 1

        global_value += 1

        return (
            outer_value,
            global_value,
        )


    return inner_function



# ============================================================
# Control flow torture
# ============================================================


def control_flow():

    for x in range(10):

        if x == 5:

            continue

        elif x == 8:

            break

        else:

            pass


    while True:

        if condition:

            return

        else:

            break



# ============================================================
# Deep indentation abuse
# ============================================================


def indentation_monster():

    if True:

        if True:

            if True:

                if True:

                    if True:

                        if True:

                            value = {
                                key:
                                    {
                                        nested:
                                            [
                                                item
                                                for item in range(10)
                                            ]
                                        for nested in range(5)
                                    }
                                for key in range(5)
                            }

                            return value



# ============================================================
# Decorator stack abuse
# ============================================================


def deco(name):

    def wrapper(fn):

        return fn

    return wrapper



@deco("one")
@deco("two")
@deco("three")
@deco(
    f"{name}:{number}"
)
@staticmethod
@classmethod
@functools.lru_cache(
    maxsize=128
)
def heavily_decorated():

    return True



# ============================================================
# try / except / else / finally combinations
# ============================================================


def try_everything():

    try:

        result = dangerous_operation()


    except ValueError as value_error:

        handle(value_error)


    except (
        TypeError,
        RuntimeError,
    ) as grouped_error:

        handle(grouped_error)


    except* Exception as exception_group:

        handle(exception_group)


    else:

        success()


    finally:

        cleanup()



# ============================================================
# Weird boolean expressions
# ============================================================


boolean_torture = (
    not (
        True
        and
        False
    )
    or
    (
        None
        is
        not
        value
    )
)



comparison_chain = (
    1
    <
    x
    <=
    y
    !=
    z
)



# ============================================================
# More slicing madness
# ============================================================


slice_torture = data[
    (
        start
    ):
    (
        stop
    ):
    (
        step
    )
]



nested_slice = (
    matrix
    [
        row_start:row_end
    ]
    [
        col_start:col_end
    ]
)



# ============================================================
# Operator precedence nightmare
# ============================================================


precedence = (
    a
    +
    b
    *
    c
    **
    d
    //
    e
    %
    f
    &
    g
    |
    h
)



# ============================================================
# End Part 3
# ============================================================

# ============================================================
# PYTHON LEXER TORTURE FILE
# PART 4/4 (VALID VERSION)
#
# FINAL BOSS
#
# Requires:
#   Python 3.14+
# ============================================================


from __future__ import annotations


import asyncio
import functools
import math
import pathlib
from collections import defaultdict
from dataclasses import dataclass
from typing import (
    Any,
    ClassVar,
    TypeAlias,
)


# ============================================================
# Keyword coverage in valid contexts
# ============================================================


def keyword_usage(
    condition: bool,
    value: Any,
):

    if condition:

        pass

    elif value is None:

        return False

    else:

        assert value is not None


    while False:

        break


    for item in range(3):

        continue


    try:

        result = value

    except Exception:

        raise

    finally:

        cleanup = True


    return result



# ============================================================
# Boolean and comparison operators
# ============================================================


boolean_expression = (
    True
    and
    False
    or
    not False
)


comparison_expression = (
    1
    <
    2
    <=
    3
    !=
    4
    ==
    4
)



membership_expression = (
    "x"
    in
    "example"
)


identity_expression = (
    None
    is
    None
)



# ============================================================
# Import edge cases
# ============================================================


import os.path as osp


from math import (
    sin,
    cos,
    tan,
)


from collections import (
    defaultdict as DD,
)



# ============================================================
# Soft keywords
# ============================================================


type UserId = int


def soft_keyword_test(value):

    match value:

        case 0:

            return "zero"

        case int(number) if number > 10:

            return "large"

        case _:

            return "other"



# ============================================================
# Async syntax
# ============================================================


async def async_function():

    await asyncio.sleep(0)


async def async_generator():

    for item in range(5):

        yield item



async def async_iteration():

    values = []

    async for item in async_generator():

        values.append(item)

    return values



# ============================================================
# Decorator stress
# ============================================================


def make_decorator(
    value: str,
):

    def decorator(
        function,
    ):

        function.metadata = value

        return function

    return decorator



@make_decorator(
    f"{__name__}:{math.pi:.5f}"
)
@functools.lru_cache(
    maxsize=128,
)
def decorated_function(
    x: int,
):

    return x * x



# ============================================================
# Classes and metaprogramming
# ============================================================


class Meta(type):

    pass



@dataclass
class ComplexClass(
    metaclass=Meta,
):

    value: int

    cache: ClassVar[
        dict[str, int]
    ] = {}


    def __post_init__(self):

        self.cache[
            str(self.value)
        ] = self.value



# ============================================================
# Slots and descriptors
# ============================================================


class Descriptor:

    def __set_name__(
        self,
        owner,
        name,
    ):

        self.name = name


    def __get__(
        self,
        instance,
        owner=None,
    ):

        return getattr(
            instance,
            self.name,
            None,
        )



class Slotted:

    __slots__ = (
        "value",
    )

    value = Descriptor()



# ============================================================
# Operator torture in valid syntax
# ============================================================


def operator_test(
    a,
    b,
):

    arithmetic = (
        a
        +
        b
        -
        a
        *
        b
        /
        b
        //
        a
        %
        b
        **
        2
    )


    bitwise = (
        (
            a
            &
            b
        )
        |
        (
            a
            ^
            b
        )
    )


    shifted = (
        a
        <<
        2
        >>
        1
    )


    matrix = (
        a
        @
        b
    )


    assignment = (
        result := arithmetic
    )


    return (
        arithmetic,
        bitwise,
        shifted,
        matrix,
        assignment,
    )



# ============================================================
# Attribute and subscription boundaries
# ============================================================


def access_torture(
    obj,
):

    return (
        obj
        .attribute
        .method(
            obj[
                0:
                10:
                2
            ]
        )
    )



# ============================================================
# String boundary tests
# ============================================================


adjacent_strings = (
    "hello"
    " "
    "world"
)


raw_strings = (
    r"\n"
    r"\t"
)


triple_strings = """
hello
{
    world
}
"""



# ============================================================
# Final f-string
# ============================================================


FINAL_F = f"""
User:
    {decorated_function(10)=}

Math:
    {math.sqrt(100):.3f}

Nested:
    {
        (
            lambda x:
                {
                    "value": x
                }
        )(42)
    }

Escapes:

    {{
        literal braces
    }}

"""



# ============================================================
# Final t-string
# ============================================================


FINAL_T = t"""
Template:

    Name:
        {__name__}

    Value:
        {math.pi}

    Literal:
        {{
            braces
        }}
"""



# ============================================================
# Indentation torture
# ============================================================


def indentation_final():

    if True:

        if True:

            for item in range(3):

                try:

                    with open(
                        "file.txt"
                    ) as handle:

                        if item:

                            result = {
                                key:
                                    [
                                        value
                                        for value in range(
                                            item
                                        )
                                    ]
                                for key in range(
                                    item
                                )
                            }

                        else:

                            result = None


                except OSError as error:

                    result = str(error)


                finally:

                    pass


    return result



# ============================================================
# EOF token
# ============================================================


FINAL_VALUE = {
    "lexer": "complete",
    "tokens":  "covered",
}