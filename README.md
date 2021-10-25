# i3Helper

collection of some i3wm related commands
### focus
easily switch to next or previous container in a workspace regardless of their floating or fullscreen properties
### snap
make windows floating and snap them to the left, right, top, or bottom of workspace.
### peek
peek through all containers in a workspace.
## Installation
you can download binaries from release page or build it yourself
First install [go](https://github.com/golang/go).

```
git clone https://github.com/smoka7/i3Helper.git
cd i3Helper
go install
```

## Usage
add commands to your i3 config

```
bindsym $mod + Up exec i3Helper focus next
bindsym $mod + Down exec i3Helper focus prev

bindsym $mod + F1 exec i3Helper snap left
bindsym $mod + F2 exec i3Helper snap top
bindsym $mod + F3 exec i3Helper snap bottom
bindsym $mod + F4 exec i3Helper snap right

# Valid time units are "ns","us","ms","s","m","h".
bindsym $mod + p exec i3Helper peek 500ms
```
