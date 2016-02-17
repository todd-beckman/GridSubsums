Subsums
=======

The goal of this project is to produce an implementation of a 2D grid m*n of numbers with two features:

- Update any single cell
- Find the sum of any subrectangle within this grid

The naive solution of this problem finds the subsum in O(mn) running time. The target running time is O(logm*logn) while retaining a reasonable (such as O(nlogn*mlogn)) space complexity.

