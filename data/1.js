// 定义一个函数，名为 findMax
function findMax(arr) {
    // 假设数组的第一个元素是最大值
    let max = arr[0];
    // 遍历数组的其余元素
    for (let i = 1; i < arr.length; i++) {
        // 如果当前元素大于最大值，更新最大值
        if (arr[i] > max) {
            max = arr[i];
        }
    }
    // 返回最大值
    return max;
}

// 测试函数
let numbers = [3, 5, 7, 9, 2, 4, 6, 8];

console.log(findMax(numbers));
console.log("Hello test world");