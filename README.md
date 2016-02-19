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
            if subquad.bounds.ContainsPoint(x, y)
                //  It is in this subquad, so update it.
                update(subquad x, y, num)
                //  The subsum has been updated, so update this one too
                quad.data = sum all 4 subquad's data
                return
            endif
        endfor
    endprocedure

####Algorithm 1 is correct

Before entering the procedure, quad.data must be the sum of all numbers bound by this quad. This is true for the first call of the procedure because the data initializes at 0, which is the sum of an empty grid. After exiting the procedure, it must be true that the target cell contains the new data and that all nodes along the quadtree used to reach this data are also updated to have new sums.

Consider the base case. If the current quad is a 1x1 cell, then it is a leaf as it cannot be split further. Thus, updating its data to the input number will trivially cause the sum of its entire subtree to be updated correctly.

Now note that because each recursive call handles half the size of the parent call, the size of the problem space is decreasing. Therefore, termination of the algorithm is guaranteed.

After the recursion has terminated, the parent of the updated node is calculated to be the sum of all 4 immediate children. Given that the 3 children that were not updated are assumed to be correct prior to running this procedure, this update is also correct. Due to the recursive nature of this algorithm, each parent's parent all the way up to the root have their sums updated correctly. Thus, the sums of all quadrants maintain correctness.

This algorithm terminates on two conditions. First, the base case which has already been shown to have correct behavior. The second, more transparent termination is if a cell is outside the bounds of the current quadrant. Because it is not contained in any of the quadrants, no recursive call is made. Thus, even faulty input will not break this algorithm.

####Algorithm 1 runs in O(log(n))

The algorithm's recursive nature can be redefined as follows, where m is the quadrant's with and n is the quadrant's height:

    T(1, 1) = 1
    T(m, n) = T(m/2, n/2) + 1

Because the width and height scale equally (halved at every iteration), they can be considered as the same variable. Effectively, whichever of m or n is greater has the more significant impact on the algorithm's behavior:

    T(n) = T(n/2) + 1

By the Master Theorem, this gives:

    T(n) = O(log(n))

###Calculating the sum of all cells in a subrectangle

Now consider the algorithm used to calculate the subsum:

    Algorithm 2
    function subsum(quad ssquad, rect SSRect)
        if quad and rect are disjoint
            return 0
        endif
        if quad is fully contained in rect
            return quad.data
        endif
        sum := 0
        parallel foreach subquad in quad: ul, ur, bl, br
            sum += subsum(subquad, rect)
        return sum
    endfunction

####Algorithm 2 is correct

The sum is initially 0 and is modified according to three conditions. If the current quad is entirely disjoint from the rectangle to sum over, then obviously none of the cells that contribute to the sum are considered. The return of 0 is appropriate at this time. Additionally, if the entire quad is contained within the rectangle, then obviously all of the cells in the quad contribute to the sum, so they are returned.

The interesting behavior is the sum made in the third case. If the quad only partially intersects with the rect, then a sum is performed over the recursive call of each subquad. Because the subquad is a quarter in size of the quad, deeper recursion will always yield calls to 1x1 rectangles which trivially must follow one of the two base cases (as a single cell cannot be just partially intersecting with the rectangle). Therefore, termination is guaranteed.

When returning from the recursive call, the sum over the four quadrants maintains correctness. That is, the sum of four 1x1 cells must be correct because each 1x1 cell reached their respective correct base cases. The sum of four correct sums, by induction, retains correctness.

####Algorithm 2 runs in O(log^2(n)) time.

With each recursive call, the work is divided into four pieces, splitting both width and height in two each time:

    T(1, 1) = 1
    T(m, n) = 4T(m/2, n/2) + 1

Again, because of the equivalent scaling as discussed earlier:

    T(n) = 4T(n/2) + 1

By Master Theorem, this gives:

    T(n) = O(n^2)

Oh no! Where is the improvement? Well, it turns out that is the work. The span of the algorithm considers the word parallel above. Because we are using Go, running programs in parallel tends to be where the power comes in. Here is the span:

    S(n) = (4 T(n/2)) / 4 + 1
    S(n) = T(n/2) + 1
    S(n) = O(log(n))

The span is O(log(n)), which is still arguably more valuable than the naive solution.

##Conclusion

So it turns out that finding the sum of n^2 numbers requires n^2 calculations. Fortunately, recursive parallelizable solutions using tree structures are quite elegant. Additionally, because of the recursive structure, it is a highly scalable solution. If the data was scattered across many disks, it would still be accessed minimally rather, while the naive solution may exhaust resources extensibly.

In the future, I may Alexify my solutions to support recursive divide-and-conquer algorithms. Twas also a great experiment for learning how to use Go.
