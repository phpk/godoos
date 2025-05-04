window.onload = function () {
    Calculator.initCache();
    Calculator.initListeners();
};

//全局计算器对象
var Calculator = (function () {
    var cal = {
        //计算器按键编码
        keyCodes: {
            0: '0',
            1: '1',
            2: '2',
            3: '3',
            4: '4',
            5: '5',
            6: '6',
            7: '7',
            8: '8',
            9: '9',
            10: '.',
            11: '±',
            12: '=',
            13: '+',
            14: '-',
            15: '*',
            16: '/',
            17: '%',
            18: '√',
            19: 'x2',
            20: '1/x',
            21: '(',
            22: ')',
            23: 'yroot',
            24: 'n!',
            25: 'Exp',
            26: '^',
            27: 'sin',
            28: 'cos',
            29: 'tan',
            30: 'powten',
            31: 'log',
            32: 'sinh',
            33: 'cosh',
            34: 'tanh',
            35: 'π',
            36: '↑',
            37: 'CE',
            38: 'C',
            39: 'Back',
            //以下是程序员型特有的按键
            40: 'A',
            41: 'B',
            42: 'C',
            43: 'D',
            44: 'E',
            45: 'F',
            46: '&',
            47: '|',
            48: '~'
        },
        //映射用于显示的操作符，比如计算时用*，而显示时x更好
        operatorFacade: {
            13: '+',
            14: '-',
            15: '×',
            16: '÷',
            17: '%',
            23: 'yroot',
            26: '^',
            46: '&',
            47: '|'
        },
        //当前计算器的类型1 --> 标准型, 2-->科学型， 3-->程序员型，默认标准型
        type: 1,
        //计算器类型前缀，用于从页面获取元素
        typePrefix: {
            1: "std-",
            2: "sci-",
            3: "pro-"
        },
        //记录每个类型的计算器的事件监听是否已经绑定,key:typpe数值，value:默认标准型是true(已加载)
        hasInited: {
            1: true,
            2: false,
            3: false
        },
        //常量
        constants: {
            //鼠标悬停时的颜色
            mouseHoverColor: "#CFCFCF",
            //计算器第一行和下面其它行的颜色是不同的，这个是第一行的背景颜色
            firstMouseOutColor: "#F2F2F2",
            //剩余各行的背景颜色
            mouseOutColor: "#E6E6E6"
        },
        cache: {
            //输入内容显示元素
            showInput: null,
            //上一步计算结果显示区域
            preStep: null,
            //显示四种进制数值的span，只在程序员型有效
            scaleSpans: null
        },
        /**
         * 获取cache.showInput的内容
         * @return String
         */
        getShowInput: function () {
            return cal.cache.showInput.innerHTML;
        },
        /**
         * 设置showInput的值
         * @param value
         */
        setShowInput: function (value) {
            cal.cache.showInput.innerHTML = value;
        },
        /**
         * 获取cache.preStep的内容
         * @return String
         */
        getPreStep: function () {
            return cal.cache.preStep.innerHTML;
        },
        setPreStep: function (value) {
            cal.cache.preStep.innerHTML = value;
        },
        //操作数栈
        operandStack: [],
        //运算符栈
        operatorStack: [],
        //上一次输入是否是二元运算符，如果是并且再次输入二元运算符，那么忽略此次输入
        isPreInputBinaryOperator: false,
        //上次按键是否是一元操作
        isPreInputUnaryOperator: false,
        //等号不可以连按
        isPreInputEquals: false,
        //如果为true，那么接下来输入的数字需要覆盖在showInput上，而不是追加
        //上一次计算的结果(=)
        preResult: 0,
        //当前使用的进制(只在程序员中有效),默认10进制(DEC)
        currentScale: 10,
        isOverride: false,
        //int校验
        intPattern: /^-?\d+$/,
        //小数校验
        floatPattern: /^-?\d+\.\d+$/,
        //科学计数法校验
        scientificPattern: /^\d+\.\d+e(\+|-)\d+$/,
        //校验16进制数字
        hexPattern: /^[0-9A-F]+$/,
        //辅助判断运算符的优先级
        operatorPriority: {
            ")": 0,
            "|": 1,
            "&": 2,
            "+": 3,
            "-": 3,
            "*": 4,
            "%": 4,
            "/": 4,
            "^": 5,
            "yroot": 5,
            "(": 6
        },
        /**
         * 初始化缓存对象(cal.cache)
         */
        initCache: function () {
            var prefix = cal.typePrefix[cal.type];
            cal.cache.showInput = document.getElementById(prefix + "show-input");
            cal.cache.preStep = document.getElementById(prefix + "pre-step");
            if (cal.type == 3) {
                cal.cache.scaleSpans = document.getElementById("pro-scales").getElementsByTagName("span");
            }
        },
        //各种事件监听函数
        listeners: {
            /**
             * 鼠标悬停在按键上的变色效果
             */
            mouseHoverListener: function (e) {
                var event = e || window.event;
                event.currentTarget.style.backgroundColor = cal.constants.mouseHoverColor;
            },
            /**
             * 鼠标从上排符号中移出的变色效果
             */
            firstMouseOutListener: function (e) {
                var event = e || window.event;
                event.currentTarget.style.backgroundColor = cal.constants.firstMouseOutColor;
            },
            /**
             * 鼠标从下排数字、符号中移出的变色效果
             */
            mouseOutListener: function (e) {
                var event = e || window.event;
                event.currentTarget.style.backgroundColor = cal.constants.mouseOutColor;
            },
            /**
             * 按键按下事件监听
             */
            keyPressListener: function (e) {
                var event = e || window.event;
                cal.handleKey(event.currentTarget.value);
            },
            /**
             * 显示/隐藏计算器类型选择栏
             */
            toggleTypeBarListener: function () {
                var bar = document.getElementById(cal.typePrefix[cal.type] + "type-bar");
                if (bar.style.display === "block") {
                    bar.style.display = "none";
                } else {
                    bar.style.display = "block";
                }
            },
            /**
             * 切换计算器类型监听器
             */
            switchTypeListener: function (e) {
                var event = e || window.event;
                cal.switchType(parseInt(event.currentTarget.value));
            },
            /**
             * 切换进制(程序员专用)
             */
            switchScaleListener: function (e) {
                var event = e || window.event;
                var scales = document.getElementById("pro-scales").getElementsByTagName("div"),
                    //此处应该使用currentTarget属性，因为target属性在绑定事件的元素有子元素的情况下会返回子元素
                    scale = parseInt(event.currentTarget.getAttribute("scale")),
                    oldScale = cal.currentScale;
                //切换选中样式
                for (var i = 0, l = scales.length; i < l; ++i) {
                    scales[i].removeAttribute("class");
                }
                event.currentTarget.setAttribute("class", "scale-active");
                var lis, btns;
                if (scale === 16) {
                    //处理上排6个16进制数字
                    cal.listeners._initFirstRowListeners();
                    if (oldScale < 10) {
                        cal.listeners._initSecondRowListeners();
                    }
                } else if (scale === 10) {
                    if (oldScale === 16) {
                        lis = document.getElementById("pro-top-symbol").getElementsByTagName("li");
                        cal.disableButtons(lis, cal.listeners.firstMouseOutListener);
                    } else {
                        cal.listeners._initSecondRowListeners();
                    }
                } else if (scale === 8) {
                    if (oldScale > 8) {
                        lis = document.getElementById("pro-top-symbol").getElementsByTagName("li");
                        cal.disableButtons(lis, cal.listeners.firstMouseOutListener);
                        //禁用8和9
                        btns = cal.getElementsByAttribute("li", "oct-disable", document.getElementById("pro-num-symbol"));
                        cal.disableButtons(btns, cal.listeners.mouseOutListener);
                    } else {
                        cal.listeners._initSecondRowListeners();
                    }
                } else if (scale === 2) {
                    if (oldScale === 16) {
                        lis = document.getElementById("pro-top-symbol").getElementsByTagName("li");
                        cal.disableButtons(lis, cal.listeners.firstMouseOutListener);
                    }
                    //禁用2-9
                    btns = cal.getElementsByAttribute("li", "bin-disable", document.getElementById("pro-num-symbol"));
                    cal.disableButtons(btns, cal.listeners.mouseOutListener);
                }
                cal.currentScale = scale;
            },
            /**
             * 初始化第一排操运算符事件监听
             * @private
             */
            _initFirstRowListeners: function () {
                var lis = document.getElementById(cal.typePrefix[cal.type] + "top-symbol").getElementsByTagName("li");
                cal.rebuildButtons(lis, cal.listeners.firstMouseOutListener);
            },
            /**
             * 初始化第二排运算符事件监听
             * @private
             */
            _initSecondRowListeners: function () {
                var lis = document.getElementById(cal.typePrefix[cal.type] + "num-symbol").getElementsByTagName("li");
                cal.rebuildButtons(lis, cal.listeners.mouseOutListener);
                if (cal.type === 3) {
                    //程序员型的小数点是禁用的
                    cal.disableButtons([document.getElementById("pro-point")], cal.listeners.mouseOutListener);
                }
            }
        },
        //初始化事件监听器
        initListeners: function () {
            var prefix = cal.typePrefix[cal.type];
            //设置上排运算符事件监听,如果是程序员型，因为默认是10进制，而上排是16进制数字，所以不需要设置事件监听
            if (cal.type < 3) {
                cal.listeners._initFirstRowListeners();
            }
            //设置下面一栏数字、四则运算事件监听
            cal.listeners._initSecondRowListeners();
            //显示/隐藏计算器类型选择侧边栏
            cal.addEvent(document.getElementById(prefix + "show-bar"), "click", cal.listeners.toggleTypeBarListener);
            //为侧边栏下的li绑定切换类型事件
            var bar = document.getElementById(prefix + "type-bar");
            lis = bar.getElementsByTagName("li");
            var li;
            for (var i = 0, l = lis.length; i < l; ++i) {
                li = lis[i];
                //非当前类型才有必要绑定事件
                if (li.className !== "active") {
                    cal.addEvent(li, "click", cal.listeners.switchTypeListener);
                }
            }
            //加载程序员型特有的
            if (cal.type === 3) {
                var scales = document.getElementById("pro-scales").getElementsByTagName("div"),
                    scale;
                for (i = 0, l = scales.length; i < l; ++i) {
                    scale = scales[i];
                    cal.addEvent(scale, "click", cal.listeners.switchScaleListener);
                }
            }
        },
        /**
         * 相应按键按下事件
         * @param value 按键的value值(即其keyCode)
         */
        handleKey: function (value) {
            var keyCode = parseInt(value);
            //如果是一个数字或者小数点，直接显示出来
            if (keyCode < 11 || (keyCode > 39 && keyCode < 46)) {
                cal.showInput(cal.keyCodes[keyCode]);
                if (cal.type === 3) {
                    //如果是程序员型，那么需要同步显示4中进制的值
                    cal.showScales(cal.getShowInput());
                }
            } else {
                switch (keyCode) {
                    //正负号
                    case 11:
                        cal.unaryOperate(function (oldValue) {
                            oldValue += "";
                            if (oldValue === "0") {
                                return [oldValue];
                            }
                            if (oldValue.charAt(0) === '-') {
                                return [oldValue.substring(1)];
                            } else {
                                return ["-" + oldValue];
                            }
                        });
                        break;
                        //开根下
                    case 18:
                        cal.unaryOperate(function (si) {
                            return [Math.sqrt(si), "sqrt"];
                        });
                        break;
                        //平方
                    case 19:
                        cal.unaryOperate(function (si) {
                            return [Math.pow(si, 2), "sqr"];
                        });
                        break;
                        //取倒数
                    case 20:
                        cal.unaryOperate(function (si) {
                            return [si === 0 ? "0不能作为被除数" : 1 / si, "1/"];
                        });
                        break;
                        //阶乘
                    case 24:
                        cal.unaryOperate(function (si) {
                            if (si < 0) {
                                si = (0 - si);
                            }
                            if (cal.isFloat(si + "")) {
                                si = Math.floor(si);
                            }
                            return [cal.fact(si), "fact"];
                        });
                        break;
                        //Exp 转为科学计数法表示
                    case 25:
                        cal.unaryOperate(function (si) {
                            return [si.toExponential(7)];
                        });
                        break;
                        //sin
                    case 27:
                        cal.unaryOperate(function (si) {
                            return [Math.sin(si), "sin"];
                        });
                        break;
                        //cos
                    case 28:
                        cal.unaryOperate(function (si) {
                            return [Math.cos(si), "cos"];
                        });
                        break;
                        //tan
                    case 29:
                        cal.unaryOperate(function (si) {
                            return [Math.tan(si), "tan"];
                        });
                        break;
                        //10的x次方
                    case 30:
                        cal.unaryOperate(function (si) {
                            return [Math.pow(10, si), "powten"];
                        });
                        break;
                        //log
                    case 31:
                        cal.unaryOperate(function (si) {
                            //js的Math.log是e的对数，Windows计算器是10的对数，此处参考Windows
                            return [Math.log10(si), "log"];
                        });
                        break;
                        //sinh(双曲正弦函数)
                    case 32:
                        cal.unaryOperate(function (si) {
                            return [Math.sinh(si), "sinh"];
                        });
                        break;
                        //cosh(双曲余弦函数)
                    case 33:
                        cal.unaryOperate(function (si) {
                            return [Math.cosh(si), "cosh"];
                        });
                        break;
                        //tanh(双曲余切函数)
                    case 34:
                        cal.unaryOperate(function (si) {
                            return [Math.tanh(si), "tanh"];
                        });
                        break;
                        //π
                    case 35:
                        cal.unaryOperate(function (si) {
                            return [Math.PI];
                        });
                        break;
                        //按位取反(~)
                    case 48:
                        cal.unaryOperate(function (si) {
                            var result = eval("~" + si);
                            //显示四种进制的数值
                            cal.showScales(result);
                            return [result];
                        });
                        break;
                        //二元运算符开始
                        //加、减、乘、除、取余，运算比较简单，直接利用eval即可求值
                    case 13:
                    case 14:
                    case 15:
                    case 16:
                    case 17:
                        //x的y次方
                    case 26:
                        //开任意次方根
                    case 23:
                        //And Or
                    case 46:
                    case 47:
                        if (cal.isPreInputBinaryOperator) {
                            break;
                        }
                        cal.isPreInputBinaryOperator = true;
                        cal.isOverride = true;
                        cal.binaryOperate(cal.keyCodes[keyCode], cal.operatorFacade[keyCode]);
                        break;
                    case 12:
                        cal.calculate();
                        break;
                        //ce
                    case 37:
                        cal.ce();
                        break;
                        //c
                    case 38:
                        cal.clear();
                        break;
                        //back
                    case 39:
                        cal.back();
                        break;
                        // (
                    case 21:
                        cal.setPreStep(cal.getPreStep() + " (");
                        cal.operatorStack.push("(");
                        break;
                        // )
                    case 22:
                        cal.rightTag();
                        break;
                        //向上箭头，把上次计算结果显示出来
                    case 36:
                        cal.setShowInput(cal.preResult);
                        break;
                }
            }
        },
        /**
         * 执行一元运算 比如取倒数、平方
         * @param operation 具体运算回调函数
         * 会向operation传递一个参数si，为用户当前的输入，同时operation函数应该返回一个数组，数组的第一个
         * 元素是计算的结果，第二个元素示例sqrt，第二个参数可选
         */
        unaryOperate: function (operation) {
            var si = cal.getShowInput(),
                result;
            if (cal.isInteger(si)) {
                result = operation(parseInt(si));
            } else if (cal.isFloat(si) || cal.isScientific(si)) {
                result = operation(parseFloat(si));
            }
            if (result != null) {
                cal.setShowInput(cal.checkLength(result[0]));
                if (result.length > 1) {
                    //显示prestep有两种情况:
                    //第一种就是这是第一次(指连续调用的第一次)调用一元函数，此时直接接在末尾即可
                    if (!cal.isPreInputUnaryOperator) {
                        cal.setPreStep(cal.getPreStep() + " " + result[1] + "(" + si + ")");
                        cal.isPreInputUnaryOperator = true;
                    } else {
                        //第二种就是这不是第一次，那么应该截取最后一个空格之后的内容进行替换
                        //比如1 + 3 + sqrt(100)，那么应该从最后一个空格后替换为此次操作的内容
                        var pi = cal.getPreStep();
                        pi = pi.substring(0, pi.lastIndexOf(" "));
                        pi += (" " + result[1] + "(" + si + ")");
                        cal.setPreStep(pi);
                    }
                }
                //一元运算结束后应该覆盖
                cal.isOverride = true;
            }
            cal.isPreInputBinaryOperator = false;
        },
        /**
         * 二元操作(+ - * / %)
         * @param operator 操作符
         * @param facade 运算符门面，用于显示在preStep中
         */
        binaryOperate: function (operator, facade) {
            //如果是程序员型，那么需要重置scalesSpan
            if (cal.type === 3) {
                cal.resetScales();
            }
            var si = cal.getShowInput(),
                pi = cal.getPreStep();
            if (cal.isNumber(si)) {
                //压操作数栈
                cal.operandStack.push(si);
                //设置preStep有三种情况:第一种上一步不是一元操作，那么需要设置si，第二种是一元操作，那么由于一元操作会把
                //函数表达式(比如sqrt(100))设置到preStep，所以不需要再次设置si
                //第三种就是如果最后一位是右括号，那么也不需要设置si
                cal.setPreStep(cal.getPreStep() + ((cal.isPreInputUnaryOperator || pi.charAt(pi.length - 1) === ")") ?
                    (" " + facade) : (" " + si + " " + facade)));
                var preOp = cal.operatorStack.pop();
                if (preOp != null) {
                    var op = cal.operatorPriority[operator],
                        pp = cal.operatorPriority[preOp];
                    //如果当前运算符优先级更高，那么只需压栈不需要计算
                    if (op > pp) {
                        cal.operatorStack.push(preOp);
                    }
                    //两者的优先级相等并且高于3(加减)，那么只需要计算一步
                    else if (op > 3 && op === pp) {
                        cal.operatorStack.push(preOp);
                        cal.travelStack(1);
                    } else {
                        cal.operatorStack.push(preOp);
                        cal.setShowInput(cal.checkLength(cal.travelStack(null, op)));
                    }
                }
                cal.operatorStack.push(operator);
            }
            cal.isPreInputUnaryOperator = false;
            cal.isPreInputEquals = false;
        },
        /**
         * 按下=时计算最终结果
         */
        calculate: function () {
            if (!cal.isPreInputEquals) {
                var si = cal.getShowInput(),
                    result;
                if (cal.isNumber(si)) {
                    cal.operandStack.push(si);
                    result = cal.checkLength(cal.travelStack());
                    cal.setShowInput(result);
                    cal.preResult = result;
                    cal.setPreStep("&nbsp;");
                    //程序员型需要把计算结果的四种进制值显示出来
                    if (cal.type === 3) {
                        cal.showScales(result);
                    }
                    cal.isOverride = true;
                }
                cal._reset();
                cal.isPreInputEquals = true;
            }
        },
        /**
         * 访问运算栈，返回计算结果
         * @param level 计算的层数，如果不指定，那么遍历整个栈
         * @param minPri(最小/截止优先级) 此参数针对下面的情况:
         * 2 + 2 X 3 X 2 ^ 2 X 2，由于最后一个运算符是X，优先级比^低，所以触发了对操作栈的遍历，但是不能全部遍历，应该遍历到第一个X停止
         * 如果不停止得到的将是错误的26 X 2 = 52，正确结果是2 + 24 X 2 = 50
         * @return Number
         * @private
         */
        travelStack: function (level, minPri) {
            var op, f, s,
                //result取操作数栈栈顶，因为防止在下列情况9 X (6 + 时出现undefined
                result = cal.operandStack[cal.operandStack.length - 1],
                l = level || cal.operatorStack.length,
                p = minPri || 0;
            for (var i = 0; i < l; ++i) {
                op = cal.operatorStack.pop();
                //遇到minPri或左括号立即停止，左括号也需要再次压入，因为只有一个右括号才能抵消一个左括号
                if (cal.operatorPriority[op] < p || op === "(") {
                    cal.operatorStack.push(op);
                    break;
                }
                s = cal.operandStack.pop();
                f = cal.operandStack.pop();
                result = cal._stackHelper(f, s, op);
                cal.operandStack.push(result);
            }
            return result;
        },
        /**
         * 输入了一个右括号
         */
        rightTag: function () {
            var si = cal.getShowInput();
            if (cal.isNumber(si)) {
                cal.setPreStep(cal.getPreStep() + (" " + si + " )"));
                cal.operandStack.push(si);
                //遍历计算操作栈，直至遇到左括号
                var op = cal.operatorStack.pop(),
                    f, s, result;
                while (op !== "(" && op != null) {
                    s = cal.operandStack.pop();
                    f = cal.operandStack.pop();
                    result = cal._stackHelper(f, s, op);
                    cal.operandStack.push(result);
                    op = cal.operatorStack.pop();
                }
                //此处应该直接把小括号的计算内容弹出，因为此结果显示在了showInput中，而再次执行二元操作时会先有一个压栈的操作，
                // 并且执行=时也是根据showInput内容计算的
                cal.setShowInput(cal.checkLength(cal.operandStack.pop()));
            }
        },
        /**
         * 辅助进行一次栈运算
         * @param f 第一个操作数
         * @param s 第二个操作数
         * @param op 运算符
         * @return 返回运算结果
         * @private
         */
        _stackHelper: function (f, s, op) {
            var result;
            if (op === "^") {
                result = Math.pow(f, s);
            } else if (op === "yroot") {
                result = Math.pow(f, 1 / s);
            }
            //+ - X / %5中操作
            else {
                //如果是程序员型，那么需要考虑进制的问题
                if (cal.type === 3) {
                    var scale = cal.currentScale,
                        fi, si;
                    if (scale === 10) {
                        result = eval(f + op + s);
                    } else if (scale === 16) {
                        fi = parseInt(f, 16);
                        si = parseInt(s, 16);
                        result = eval(fi + op + si).toString(16);
                    } else if (scale === 8) {
                        fi = parseInt(f, 8);
                        si = parseInt(s, 8);
                        result = eval(fi + op + si).toString(8);
                    } else {
                        fi = parseInt(f, 2);
                        si = parseInt(s, 2);
                        result = eval(fi + op + si).toString(2);
                    }
                } else {
                    result = eval(f + op + s);
                }
            }
            return result;
        },
        /**
         * 确保结果长度不大于13,如果超出，以科学计数法形式显示(小数点后7位)
         * @param value 需要检查的结果
         */
        checkLength: function (value) {
            var valueStr = value + "";
            if (cal.isFloat(valueStr)) {
                valueStr = valueStr.replace(/0+$/, "");
            }
            return valueStr.length > 12 ? value.toExponential(7) : valueStr;
        },
        //CE
        ce: function () {
            cal.setShowInput("0");
            if (cal.type === 3) {
                cal.resetScales();
            }
        },
        //C
        clear: function () {
            cal.setShowInput("0");
            cal.setPreStep("&nbsp;");
            cal._reset();
            if (cal.type === 3) {
                cal.resetScales();
            }
        },
        /**
         * 清空四个进制的值
         * @private
         */
        resetScales: function () {
            for (var i = 0; i < 4; i++) {
                cal.cache.scaleSpans[i].innerHTML = "0";
            }
        },
        back: function () {
            var oldValue = cal.cache.showInput.innerText;
            cal.setShowInput(oldValue.length < 2 ? "0" : oldValue.substring(0, oldValue.length - 1));
        },
        /**
         * 当计算器类型是程序员时，需要同步显示四种进制的值
         * @param num 需要显示的数字
         */
        showScales: function (num) {
            var result = cal.calculateScales(num),
                spans = cal.cache.scaleSpans;
            for (var i = 0; i < 4; ++i) {
                spans[i].innerHTML = result[i];
            }
        },
        /**
         * 根据当前进制分别计算出四种进制的值
         * @param num 需要计算的值
         * @return Array 共4个元素，依次为16、10、8、2进制的值
         */
        calculateScales: function (num) {
            var scale = cal.currentScale,
                result = [],
                i;
            if (scale === 10) {
                i = parseInt(num);
                result[0] = i.toString(16);
                result[1] = i;
                result[2] = i.toString(8);
                result[3] = i.toString(2);
            } else if (scale === 16) {
                //先转成10进制，然后再转成其它进制
                i = parseInt(num, 16);
                result[0] = num;
                result[1] = i;
                result[2] = i.toString(8);
                result[3] = i.toString(2);
            } else if (scale === 8) {
                i = parseInt(num, 8);
                result[0] = i.toString(16);
                result[1] = i;
                result[2] = num;
                result[3] = i.toString(2);
            } else {
                i = parseInt(num, 2);
                result[0] = i.toString(16);
                result[1] = i;
                result[2] = i.toString(8);
                result[3] = num;
            }
            return result;
        },
        /**
         * 校验字符串是否是数字
         * @param str
         * @return 是返回true
         */
        isNumber: function (str) {
            return cal.isInteger(str) || cal.isFloat(str) || cal.isScientific(str) || cal.isHex(str);
        },
        /**
         * 校验是否是整数
         * @param str
         */
        isInteger: function (str) {
            return str.match(cal.intPattern);
        },
        /**
         * 校验是否是小数
         * @param str
         */
        isFloat: function (str) {
            return str.match(cal.floatPattern);
        },
        /**
         * 是否是科学计数法
         * @param str
         */
        isScientific: function (str) {
            return str.match(cal.scientificPattern);
        },
        /**
         * 是否是16进制数字
         * @param str
         */
        isHex: function (str) {
            return str.match(cal.hexPattern);
        },
        /**
         * 显示输入的内容
         * 用于相应数字/小数点按键
         * @param value 按键的内容，不是keyCode
         */
        showInput: function (value) {
            var oldValue = cal.getShowInput();
            var newValue = oldValue;
            if (cal.isOverride) {
                //既然是覆盖，那么如果直接输入.那么肯定是0.x
                if (value === ".") {
                    newValue = "0.";
                } else {
                    newValue = value;
                }
            } else if (oldValue.length < 13) {
                if (oldValue === "0") {
                    if (value === ".") {
                        newValue = "0.";
                    } else {
                        newValue = value;
                    }
                } else {
                    newValue += value;
                }
            }
            cal.setShowInput(newValue);
            cal.isOverride = false;
            cal.isPreInputBinaryOperator = false;
            cal.isPreInputUnaryOperator = false;
            cal.isPreInputEquals = false;
        },
        /**
         * 切换计算器类型
         * @param type int 要切换到的类型
         */
        switchType: function (type) {
            //关闭选择栏
            var oldPrefix = cal.typePrefix[cal.type];
            document.getElementById(oldPrefix + "type-bar").style.display = "none";
            //切换面板
            document.getElementById(oldPrefix + "main").style.display = "none";
            document.getElementById(cal.typePrefix[type] + "main").style.display = "block";
            cal.type = type;
            if (!cal.hasInited[type]) {
                cal.initListeners();
                cal.hasInited[type] = true;
            }
            cal.initCache();
            cal._reset();
        },
        /**
         * 重置各个标志变量以及操作栈
         * @private
         */
        _reset: function () {
            cal.operandStack = [];
            cal.operatorStack = [];
            cal.isPreInputBinaryOperator = false;
            cal.isPreInputUnaryOperator = false;
            cal.isPreInputEquals = false;
        },
        /**
         * 工具方法，为element添加事件处理函数
         * @param element 需要添加事件的dom元素
         * @param name name事件名称(不含on)
         * @param handler 事件处理函数
         */
        addEvent: function (element, name, handler) {
            if (window.addEventListener) {
                element.addEventListener(name, handler);
            } else if (window.attachEvent) {
                element.attachEvent("on" + name, handler);
            }
        },
        /**
         * 工具方法，为element移除特定的事件监听
         * @param element 需要移除事件监听的dom元素
         * @param name 事件名，没有"on"
         * @param handler 需要移除的处理函数
         */
        removeEvent: function (element, name, handler) {
            if (window.removeEventListener) {
                element.removeEventListener(name, handler);
            } else if (window.detachEvent) {
                element.detachEvent("on" + name, handler);
            }
        },
        /**
         * 根据元素的属性进行查找，只要存在此属性即可
         * @param tag 目标标签名
         * @param attr
         * @param root 开始查找的起始节点，可选，默认document
         */
        getElementsByAttribute: function (tag, attr, root) {
            var parent = root || document,
                result = [];
            var arr = parent.getElementsByTagName(tag),
                a;
            for (var i = 0, l = arr.length; i < l; ++i) {
                a = arr[i];
                if (a.getAttribute(attr) != null) {
                    //这个写法...
                    result[result.length] = a;
                }
            }
            return result;
        },
        /**
         * 阶乘
         * @param n 操作数 int
         * @return
         */
        fact: (function () {
            //缓存
            var cache = [1];

            function factorial(n) {
                var result = cache[n - 1];
                if (result == null) {
                    result = 1;
                    for (var i = 1; i <= n; ++i) {
                        result *= i;
                    }
                    cache[n - 1] = result;
                }
                return result;
            }
            return factorial;
        })(),
        /**
         * 禁用按键，只有程序员型才会用到
         * @param lis 按钮集合
         * @param mouseOutListener function 鼠标移出时采用哪个监听函数，取决于按钮的位置(上排/下排)
         */
        disableButtons: function (lis, mouseOutListener) {
            var li;
            for (var i = 0, l = lis.length; i < l; ++i) {
                li = lis[i];
                li.setAttribute("class", "disable-btn");
                cal.removeEvent(li, "click", cal.listeners.keyPressListener);
                cal.removeEvent(li, "mouseout", mouseOutListener);
                cal.removeEvent(li, "mouseover", cal.listeners.mouseHoverListener);
            }
        },
        /**
         * 重新设置按键
         * @param lis 按钮集合
         * @param mouseOutListener function 鼠标移出时采用哪个监听函数，取决于按钮的位置(上排/下排)
         */
        rebuildButtons: function (lis, mouseOutListener) {
            var li;
            for (var i = 0, l = lis.length; i < l; ++i) {
                li = lis[i];
                li.removeAttribute("class");
                cal.addEvent(li, "click", cal.listeners.keyPressListener);
                cal.addEvent(li, "mouseout", mouseOutListener);
                cal.addEvent(li, "mouseover", cal.listeners.mouseHoverListener);
            }
        }
    };
    return cal;
})();