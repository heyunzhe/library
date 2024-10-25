// 获取模态窗口
var modal = document.getElementById("bookshelfModal");

// 获取打开模态窗口的按钮
var btn = document.getElementById("myBookshelf");

// 获取关闭按钮
var span = document.getElementsByClassName("close")[0];

let bookCount = 0;
let borrowedBooksArray = []; // 用于存储借阅的书籍元素

// 当用户点击按钮时，打开模态窗口
btn.onclick = function() {
    if (localStorage.getItem('isLoggedIn') === 'true') {
        modal.style.display = "block";
    } else {
        window.location.href = '/login';
    }
};

// 当用户点击 (x), 关闭模态窗口
span.onclick = function() {
    modal.style.display = "none";
}

// 当用户点击模态窗口外部时，关闭它
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const authButton = document.getElementById('authButton');
    
    // 检查登录状态
    if (localStorage.getItem('isLoggedIn') === 'true') {
        authButton.textContent = '个人中心';
        authButton.href = '/user/library';
    } else {
        authButton.textContent = '登录';
        authButton.href = '/login';
    }

    // 使用事件委托来处理动态添加的 "加入借书架" 按钮
    document.querySelector('.box13').addEventListener('click', function(event) {
        if (event.target.matches('#lend_book')) {
            event.preventDefault();
            
            if (localStorage.getItem('isLoggedIn') === 'true') {
                const bookInfo = event.target.closest('[class^="box13-"]');
                const bookTitle = bookInfo.querySelector('h3').textContent;
                const isbn = bookInfo.querySelector('.ISBN').textContent;
                
                addToBookshelf(bookTitle, isbn);
            } else {
                window.location.href = '/login';
            }
        }
    });

    // 使用事件委托，将监听器添加到一个始终存在的父元素上
    document.getElementById('borrowedBooks').addEventListener('submit', function(event) {
        if (event.target.matches('#lendform')) {
            event.preventDefault(); // 阻止表单的默认提交

            // 弹出确认对话框
            if (confirm("您确定要借阅这本书吗？")) {
                var formData = new FormData(event.target);
                var bookElement = event.target.closest('div');
                
                fetch('/lend/book', {
                    method: 'POST',
                    body: formData,
                })
                .then(response => {
                    if (response.ok) {
                        alert("已成功借出");
                        removeFromBookshelf(bookElement);
                    } else if (response.status === 400) {
                        alert("日期设置错误");
                    } else if (response.status === 403) {
                        alert("已借阅此书或已达借阅上限");
                    } else if (response.status === 404) {
                        alert("此书已被借完");
                    } else {
                        alert("服务器错误");
                    }
                })
                .catch(error => {
                    console.error('请求错误:', error);
                });
            }
        }
    });
});

function addToBookshelf(bookTitle, isbn) {
    const borrowedBooks = document.getElementById('borrowedBooks');
    const bookElement = document.createElement('div');
    bookElement.innerHTML = `
        <form action="/lend/book" method="post" id="lendform">    
        <p>${bookTitle}</p>
        <p>ISBN：${isbn}</p>
        预计归还日期：<input type="date" name="exp_return_date" required>
        <br>
        <input type="hidden" name="isbn" value="${isbn}">
        <br>
        <button type="submit">确认借阅</button>
        </form>
        <button class="delete-button">删除</button>
    `;
    
    borrowedBooks.appendChild(bookElement);
    borrowedBooksArray.push(bookElement);

    // 增加书籍计数
    bookCount++;
    // 更新借书架上的书籍数量显示
    document.getElementById('bookCount').textContent = `（${bookCount}）`;

    // 为删除按钮添加事件监听器
    bookElement.querySelector('.delete-button').addEventListener('click', function() {
        removeFromBookshelf(bookElement);
    });
}

function removeFromBookshelf(bookElement) {
    const borrowedBooks = document.getElementById('borrowedBooks');
    borrowedBooks.removeChild(bookElement);
    borrowedBooksArray = borrowedBooksArray.filter(item => item !== bookElement);
    
    // 减少书籍计数
    bookCount--;
    // 更新借书架上的书籍数量显示
    document.getElementById('bookCount').textContent = `（${bookCount}）`;
}