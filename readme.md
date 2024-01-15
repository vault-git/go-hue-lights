# Project
GO cli tool that enables you to control Phillips smart lights through the HUE Bridge.

## Example

```
Usage of ./go-hue-light:
  -br float
        Controls the brightness of the given light. [0 - 100] (default 100)
  -colorx float
        Controls the X Coordinate in the color diagram. [0.0 - 1.0]
  -colory float
        Controls the Y Coordinate in the color diagram. [0.0 - 1.0]
  -light string
        Name of the light to control
```

## Todo
- command to get detailed light status

## Bugs
- when ommiting the color flags, the color changes to the default values
