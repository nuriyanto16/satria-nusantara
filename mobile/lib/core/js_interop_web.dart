import 'dart:js' as js;

dynamic get jsContext => js.context;
dynamic allowInterop(Function fn) => js.allowInterop(fn);
