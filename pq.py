def solution(A):
    N = len(A)
    if N < 2:
        return 0

    dp = [[0] * 4 for _ in range(N)]  # dp[i][k]: max sum at index i using k tiles

    for i in range(1, N):
        pair_sum = A[i - 1] + A[i]

        for k in range(1, 4):  # 1 to 3 tiles
            # case 1: Not using tile at (i-1, i)
            dp[i][k] = max(dp[i][k], dp[i - 1][k])

            # case 2: Using tile at (i-1, i)
            if i >= 2:
                dp[i][k] = max(dp[i][k], dp[i - 2][k - 1] + pair_sum)
            else:
                # If i == 1, we can only look back to i == -1, so just add pair
                dp[i][k] = max(dp[i][k], pair_sum)

        dp[i][0] = dp[i - 1][0]

    return max(dp[N - 1])


if __name__ == "__main__":
    print(solution([2, 3, 5, 2, 3, 4, 6, 4, 1]))     # ➜ 25
    print(solution([1, 5, 3, 2, 6, 6, 10, 4, 7, 2, 1]))  # ➜ 35
    print(solution([1, 2, 3, 3, 2]))               # ➜ 10
    print(solution([5, 10, 3]))                   # ➜ 15
