Output assets/como-usar.gif

# Requisitos
Require echo
Require go


# Config visual
Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 800
Set LetterSpacing 1
Set LineHeight 1
Set Theme "Catppuccin Frappe"
Set Padding 10
Set Framerate 60
Set LoopOffset 3


# Execucao
Hide
Type "clear; go run ~/estudo/go/market2csv main.go" Enter
Show
Sleep 3

Screenshot assets/pesquisa.png
Type@250ms "celular xiaomi 5g" Enter
Sleep 5
Screenshot assets/quantidade.png

Type@200ms "15" Enter
Sleep 18
Screenshot assets/final.png