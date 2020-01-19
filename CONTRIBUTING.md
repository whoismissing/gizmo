# Contribution

Additional service checks can be implemented by writing a function in [gizmo/check](https://github.com/whoismissing/gizmo/tree/master/check)

Then, a corresponding structure and CheckHealth function implementing the ServiceCheck interface should be written in [gizmo/structs](https://github.com/whoismissing/gizmo/tree/master/structs)

Finally, a case statement must be added to the LoadFromServiceType function in [gizmo/structs](https://github.com/whoismissing/gizmo/tree/master/structs)
