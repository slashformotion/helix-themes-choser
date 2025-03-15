
Set Padding 0
Set Width 1920
Set Height 1080
Set FontSize 30

Require hx
Type "hx {{ file }} -c config.toml"
Enter
Sleep 500ms
Type ":theme {{ theme }}"
Sleep 10ms
Enter
Sleep 100ms
Screenshot {{ out }}.png
Sleep 300ms
Type ":q"
