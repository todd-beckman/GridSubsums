#Subsums

This project's implementation is complete, though thorough testing has not been done.


The goal of this project is to produce an implementation of a 2D grid mxn of numbers with two features:

- Update any single cell
- Find the sum of any subrectangle within this grid

The naive solution of this problem finds the subsum in O(mn) running time. The target running time is O(log(mn)) while retaining a reasonable space complexity. A third power polynomial or worse is restricted.


##Theory

Here we prove the correctness and running time/space of the algorithms and discuss the data structures involved.

###Quadtree

The data structure of choice is a quadtree. Each node along this tree (hereafter referred to as "quad") is responsible for a quadrant which contains one quarter the area of the parent node. The "ul" node is, for example, responsible for all data in the upper lefthand corner. The leaf nodes are 1x1 rectangles whose data are the individual cell values, while each parent (and their parents and so on) has data equal to the sum of all leaves in its subtree.

This implementation increases the runtime of the update function to O(log(mn)), but it decreases the runtime of the subsum function to O(log(mn)). The space required to maintain this data structure is O(mn x log(mn)).

###Updating a Grid Cell

Consider first the algorithm used for updating a single cell:

    Algorithm 1
    procedure update(quad ssquad, x, y, num int)
        //  base case: looking at a single cell
        if quad.bounds is a 1x1 rectangle at (x,y)
            quad.data = num
            return
        end if
        for each subquad in quad: ul, ur, bl, br
            //  construct null subquadrants
            if subquad == nil
                subquad = new ssquad with bounds splitting quad in 4 parts
            endif
            if subquad.bounds.ContainsPoint(x, y)
                //  It is in this subquad, so update it.
                update(subquad x, y, num)
                //  The subsum has been updated, so update this one too
                quad.data = sum all 4 subquad's data
                return
            endif
        endfor
    endprocedure

###Algorithm 1 is correct.

Before entering the procedure, quad.data must be the sum of all numbers bound by this quad. This is true for the first call of the procedure because the data initializes at 0, which is the sum of an empty grid. After exiting the procedure, it must be true that the target cell contains the new data and that all nodes along the quadtree used to reach this data are also updated to have new sums.

Consider the base case. If the current quad is a 1x1 cell, then it is a leaf as it cannot be split further. Thus, updating its data to the input number will trivially cause the sum of its entire subtree to be updated correctly.

After the recursion has terminated, the parent of the updated node is calculated to be the sum of all 4 immediate children. Given that the 3 children that were not updated are assumed to be correct prior to running this procedure, this update is also correct. Due to the recursive nature of this algorithm, each parent's parent all the way up to the root have their sums updated correctly. Thus, the sums of all quadrants maintain correctness.

This algorithm terminates on two conditions. First, the base case which has already been shown to have correct behavior. The second, more transparent termination is if a cell is outside the bounds of the current quadrant. Because it is not contained in any of the quadrants, no recursive call is made. Thus, even faulty input will not break this algorithm.

###Algorithm 1 runs in O(log(mn))

The algorithm's recursive nature can be redefined as follows, where m is the quadrant's with and n is the quadrant's height:

    T(1, 1) = 1
    T(m, n) = T(m/2, n/2) + 4

By the Master Theorem, this gives:

    T(m, n) = O(log(m) log(n))
    T(m, n) = O(log(mn))

