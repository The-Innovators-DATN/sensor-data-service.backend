def solve(x, y):
    if len(x) != len(y):
        print("NO")  # Mismatched inputs
        return

    n = len(x)
    if n < 3:
        print("YES")  # Any 2 or fewer points can always lie on a square perimeter
        return

    min_x, max_x = min(x), max(x)
    min_y, max_y = min(y), max(y)
    
    # Determine the target square length
    square_length = max(max_x - min_x, max_y - min_y)
    
    # Adjust x range to match square length
    if (max_x - min_x) < square_length:
        min_x = min(max_x - square_length, min_x)
        max_x = max(min_x + square_length, max_x)
    
    # Adjust y range to match square length
    if (max_y - min_y) < square_length:
        min_y = min(max_y - square_length, min_y)
        max_y = max(min_y + square_length, max_y)
    
    x_range = (min_x, max_x)
    y_range = (min_y, max_y)
    
    for i in range(n):
        if not (x_range[0] <= x[i] <= x_range[1] and y_range[0] <= y[i] <= y_range[1]):
            print("NO")


if __name__ == "__main__":
    # x = [-2, -1, 2, 2, 1, -2]
    # y = [1, 1, 1, -1, -1, -1]
    # OUtput: No, cause it is rectangle
    x = [-2, -1, 2, 2, 1, -2]
    y = [1, 1, 1, -1, -1, -1]
    # x = [-1, 1 , -2]
    # y = [2, 1, 0]
    # Output: YES, cause it is square, and to know the 4 edge, we must use the square_length
    solve(x, y)