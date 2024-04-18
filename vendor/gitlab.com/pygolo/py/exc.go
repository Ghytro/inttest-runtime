/*
 * Copyright 2022, Pygolo Project contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package py

// #include "pygolo.h"
import "C"

var Exc_BaseException = Object{C.PyExc_BaseException}
var Exc_Exception = Object{C.PyExc_Exception}

// Standard Python exceptions.
//
// C API: https://docs.python.org/3/c-api/exceptions.html#standard-exceptions
var Exc_ArithmeticError = Object{C.PyExc_ArithmeticError}
var Exc_AssertionError = Object{C.PyExc_AssertionError}
var Exc_AttributeError = Object{C.PyExc_AttributeError}
var Exc_BlockingIOError = Object{C.PyExc_BlockingIOError}
var Exc_BrokenPipeError = Object{C.PyExc_BrokenPipeError}
var Exc_BufferError = Object{C.PyExc_BufferError}
var Exc_ChildProcessError = Object{C.PyExc_ChildProcessError}
var Exc_ConnectionAbortedError = Object{C.PyExc_ConnectionAbortedError}
var Exc_ConnectionError = Object{C.PyExc_ConnectionError}
var Exc_ConnectionRefusedError = Object{C.PyExc_ConnectionRefusedError}
var Exc_ConnectionResetError = Object{C.PyExc_ConnectionResetError}
var Exc_EOFError = Object{C.PyExc_EOFError}
var Exc_FileExistsError = Object{C.PyExc_FileExistsError}
var Exc_FileNotFoundError = Object{C.PyExc_FileNotFoundError}
var Exc_FloatingPointError = Object{C.PyExc_FloatingPointError}
var Exc_GeneratorExit = Object{C.PyExc_GeneratorExit}
var Exc_ImportError = Object{C.PyExc_ImportError}
var Exc_IndentationError = Object{C.PyExc_IndentationError}
var Exc_IndexError = Object{C.PyExc_IndexError}
var Exc_InterruptedError = Object{C.PyExc_InterruptedError}
var Exc_IsADirectoryError = Object{C.PyExc_IsADirectoryError}
var Exc_KeyError = Object{C.PyExc_KeyError}
var Exc_KeyboardInterrupt = Object{C.PyExc_KeyboardInterrupt}
var Exc_LookupError = Object{C.PyExc_LookupError}
var Exc_MemoryError = Object{C.PyExc_MemoryError}
var Exc_ModuleNotFoundError = Object{C.PyExc_ModuleNotFoundError}
var Exc_NameError = Object{C.PyExc_NameError}
var Exc_NotADirectoryError = Object{C.PyExc_NotADirectoryError}
var Exc_NotImplementedError = Object{C.PyExc_NotImplementedError}
var Exc_OSError = Object{C.PyExc_OSError}
var Exc_OverflowError = Object{C.PyExc_OverflowError}
var Exc_PermissionError = Object{C.PyExc_PermissionError}
var Exc_ProcessLookupError = Object{C.PyExc_ProcessLookupError}
var Exc_RecursionError = Object{C.PyExc_RecursionError}
var Exc_ReferenceError = Object{C.PyExc_ReferenceError}
var Exc_RuntimeError = Object{C.PyExc_RuntimeError}
var Exc_StopAsyncIteration = Object{C.PyExc_StopAsyncIteration}
var Exc_StopIteration = Object{C.PyExc_StopIteration}
var Exc_SyntaxError = Object{C.PyExc_SyntaxError}
var Exc_SystemError = Object{C.PyExc_SystemError}
var Exc_SystemExit = Object{C.PyExc_SystemExit}
var Exc_TabError = Object{C.PyExc_TabError}
var Exc_TimeoutError = Object{C.PyExc_TimeoutError}
var Exc_TypeError = Object{C.PyExc_TypeError}
var Exc_UnboundLocalError = Object{C.PyExc_UnboundLocalError}
var Exc_UnicodeDecodeError = Object{C.PyExc_UnicodeDecodeError}
var Exc_UnicodeEncodeError = Object{C.PyExc_UnicodeEncodeError}
var Exc_UnicodeError = Object{C.PyExc_UnicodeError}
var Exc_UnicodeTranslateError = Object{C.PyExc_UnicodeTranslateError}
var Exc_ValueError = Object{C.PyExc_ValueError}
var Exc_ZeroDivisionError = Object{C.PyExc_ZeroDivisionError}
