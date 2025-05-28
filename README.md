# align
A lightweight, framework-agnostic UI layout library for Go

## Design Philosophy

Most UI frameworks take an "all-in-one" approach, providing automatic resizing, complex layout managers, and hierarchical component trees. While powerful, these systems often introduce significant complexity for what should be simple positioning tasks. The `align` library takes a radically different approach: **it focuses exclusively on positioning relationships, deliberately omitting automatic sizing**.

Instead of managing parent-child hierarchies with cascading layout effects, `align` uses linked lists to group related UI elements. Rather than implementing complex auto-resize logic that can trigger unpredictable layout recalculations, it requires manual size specification. This intentional limitation transforms UI layout from a complex, often unpredictable system into a simple, transparent tool.

The key insight is that in many applications—especially games—element sizes are often predetermined or calculated procedurally. What developers really need is an elegant way to express spatial relationships: "place button B below button A," "center this dialog," or "align this group to the bottom-right corner." Traditional layout systems bury these simple concepts under layers of box models, constraints, and automatic behaviors.

By embracing manual size management, `align` eliminates entire categories of bugs related to unexpected resizing, constraint conflicts, and layout cascades. The resulting code is more predictable, easier to debug, and significantly more performant. What you lose in automation, you gain in simplicity, transparency, and control.

This "less is more" philosophy makes `align` particularly well-suited for game development, where UI elements typically have fixed dimensions and the primary challenge is positioning them relative to each other and the screen boundaries.
