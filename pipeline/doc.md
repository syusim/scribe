For now this package just exists for me to lay out some baseline thoughts about
how to structure a file watcher/auto-recompilation setup.

My thoughts right now are poisoned due to having thought about Timely Dataflow
a lot recently but I think a simplified version of that could be good.  There's
a base "file watcher" vertex which emits file system events to other worker
vertices, like one that takes all the posts and compiles a table of contents
from them, and one that takes each post individually and compiles it to html.

Not sure if more generality than that is required?
