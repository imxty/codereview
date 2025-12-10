sum_num = 0
i = 1
while i <= 5:
    sum_num += i
# 缺少i += 1，i永远等于1，循环无法终止


a = 10
b = "20"
print(a + int(b))  # 转换为相同类型（数字）
# 或
print(str(a) + b)  # 转换为相同类型（字符串）


# 错误代码
print('Hello, Python")  # 单引号开头，双引号结尾

# 错误代码：想判断年龄是否在18-60之间，结果逻辑相反
age = 25
if age < 18 or age > 60:
    print("成年且未退休")
else:
    print("未成年或已退休")