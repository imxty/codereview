while i <= 5:
    sum_num += i
    i += 1
# 缺少i += 1，i永远等于1，循环无法终止


# 错误代码
print("Hello, Python")  # 单引号开头，双引号结尾
      
# 错误代码：想判断年龄是否在18-60之间，结果逻辑相反
age = 25
if age >= 18 and age <60:
    print("成年且未退休")
else:
    print("未成年或已退休")