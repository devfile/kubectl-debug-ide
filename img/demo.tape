Output demo.gif
Require kubectl
Require bat
Require cat

Set Shell "zsh"
Set FontSize 23
Set Width 1400
Set Height 550
# Set Padding 5
# Set FontFamily "Iosevka Term"
# Set Theme "AdventureTime"
# Set Theme "Aurora"

# Hide
# Type 'alias cat="bat -pP"'
# Enter
# Type clear
# Enter
# Sleep 500ms
# Show

Type "kubectl get pods"
Sleep 500ms
Enter
Sleep 2s
Enter
Enter
Type "kubectl debug-ide outyet \"
Enter
Type "  --image ghcr.io/l0rd/outyet-dev:latest \"
Enter
Type "  --copy-to outyet-debug \"
Enter
Type "  --share-processes \"
Enter
Type "  --git-repository https://github.com/l0rd/outyet.git"
Sleep 500ms
Enter
Sleep 15s
