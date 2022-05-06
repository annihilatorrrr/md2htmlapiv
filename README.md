# gotg_md2html api

This is the api version of go md to html; just deploy it and make own request:

A working version is here:
[Click!](https://md2htmlapi.herokuapp.com/)
```
# python 3
import requests as r
e = r.post("https://md2htmlapi.herokuapp.com/md2htmlbv2/", data={"rtext": "some text from user"})
print(e.json())
```
Getting one button:
> For more than one use for loop
```
import requests as r
e = r.post("https://https://md2htmlapi.herokuapp.com/md2htmlbv2/", data={"rtext": m.reply_to_message.text.markdown})
J = list((e.json().get("button").replace("[", "").replace("]", "").replace("{", "").replace("}", "")).split(" "))
print(J[len(J) - 1:][0]) # button same_line
print(J[len(J) - 2:][0]) # button link
K = J[:-2]
print(" ".join([K[s] for s, b in enumerate(K)])) # button text
```
[![DeepSource](https://deepsource.io/gh/annihilatorrrr/md2htmlapiv.svg/?label=active+issues&show_trend=true&token=bkO4NXDBYGKNjgh-gC7RQ0ah)](https://deepsource.io/gh/annihilatorrrr/md2htmlapiv/?ref=repository-badge)
# Credits:
[Divkix](https://github.com/Divkix); Formatting fixer.

[AmarnathCJD](https://github.com/AmarnathCJD); Every push-ups to fix issues.

[Anony](https://github.com/anonyindian); The idea giver and what to do and how?
