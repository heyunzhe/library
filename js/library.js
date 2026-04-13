document.addEventListener('DOMContentLoaded', function() {
    const toggleButton = document.getElementById('togglebutton');
    const searchDiv = document.getElementById('search_div');
    const bookSearch = document.getElementById('book_search');
    const webSearch = document.getElementById('web_search');
    const searchSelect = document.getElementById('search_select');
    const searchInput = document.querySelector('.search');
    const searchButton = document.querySelector('.box9');
    const box6_1 = document.querySelector('.box6-1');

    let isSearchDivVisible = false;

    // 点击搜索图书资源时切换显示状态
    toggleButton.addEventListener('click', function() {
        isSearchDivVisible = !isSearchDivVisible; // 切换状态
        searchDiv.style.display = isSearchDivVisible ? 'block' : 'none';

        if (isSearchDivVisible) {
            bookSearch.click(); // 默认选择图书搜索
        }
    });

    
    // 点击图书搜索
    bookSearch.addEventListener('click', function() {
        // 移除站内搜索的活跃状态
        webSearch.classList.remove('active');
        // 添加图书搜索的活跃状态
        bookSearch.classList.add('active');
        // 更新 box6-1 的文字
        box6_1.textContent = "图书搜索";
        // 显示 select 标签和输入框
        searchSelect.style.display = 'block';
        searchInput.style.marginLeft = '0'; // 还原输入框位置
        searchButton.style.marginLeft = '0'; // 还原按钮位置
    });

    // 点击站内搜索
    webSearch.addEventListener('click', function() {
        // 移除图书搜索的活跃状态
        bookSearch.classList.remove('active');
        // 添加站内搜索的活跃状态
        webSearch.classList.add('active');
        // 更新 box6-1 的文字
        box6_1.textContent = "站内搜索";
        // 隐藏 select 标签
        searchSelect.style.display = 'none';
        searchInput.style.marginLeft = '-120px'; // 输入框前移
        searchButton.style.marginLeft = '-120px'; // 按钮前移
    });
});


const images = [
            'url("../images/home.jpg")',
            'url("../images/home1.jpg")',
            'url("../images/home2.jpg")',
            'url("../images/home3.jpg")'

        ];
        let currentIndex = 0;

        function changeBackground() {
            currentIndex = (currentIndex + 1) % images.length; // 循环索引
            document.getElementById('fbox').style.backgroundImage = images[currentIndex];
        }

        // 每五秒自动切换背景图片
        setInterval(changeBackground, 5000);


// 定义书籍数据
const recommendedBooks = [
    { title: "活着", author: "余华", image: "../images/image1.jpg" },
    { title: "提问的艺术", author: "[美] 特里.费德姆", image: "../images/image2.jpg" },
    { title: "狂人日记", author: "鲁迅", image: "../images/image3.jpg" },
    { title: "钢铁是怎样炼成的", author: "（苏）尼·奥斯特洛夫斯基", image: "../images/image4.jpg" },
    { title: "骆驼祥子", author: "老舍", image: "../images/image5.jpg" }
];

const newBooks = [
    { title: "海边的卡夫卡", author: "村上春树", image: "../images/image6.jpg" },
    { title: "百年孤独", author: "加西亚·马尔克斯", image: "../images/image7.jpg" },
    { title: "追风筝的人", author: "卡勒德·胡赛尼", image: "../images/image8.jpg" },
    { title: "动物农场", author: "乔治·奥威尔", image: "../images/image9.jpg" },
    { title: "挪威的森林", author: "村上春树", image: "../images/image10.jpg" }
];

// 获取元素
const recBookDiv = document.querySelector('.rec_book');
const newBookDiv = document.querySelector('.new_book');
const boxes = document.querySelectorAll('.box16, .box17, .box18, .box19, .box20');

let isRecommended = true; // 默认状态为好书推荐

function displayBooks(books) {
    boxes.forEach((box, index) => {
        const book = books[index];
        box.querySelector('img').src = book.image;
        box.querySelector('.title').textContent = book.title;
        box.querySelector('.author').textContent = book.author;
    });
}

// 初始化显示好书推荐
displayBooks(recommendedBooks);

function toggleBooks() {
    if (isRecommended) {
        // 切换到新书
        recBookDiv.style.backgroundColor = 'rgb(0,0,0)'; 
        recBookDiv.style.color = 'rgb(255,255,255)';
        newBookDiv.style.backgroundColor = 'rgb(46, 200, 131)'; 
        displayBooks(newBooks); // 显示新书
    } else {
        // 切换到好书推荐
        newBookDiv.style.backgroundColor = 'rgb(0, 0, 0)'; 
        newBookDiv.style.color = 'rgb(255,255,255)'; 
        recBookDiv.style.backgroundColor = 'rgb(46, 200, 131)'; 
        displayBooks(recommendedBooks); // 显示好书推荐
    }
    isRecommended = !isRecommended; // 切换状态

    // 禁用当前按钮
    recBookDiv.style.pointerEvents = isRecommended ? 'none' : 'auto';
    newBookDiv.style.pointerEvents = isRecommended ? 'auto' : 'none';
}

// 添加点击事件
newBookDiv.addEventListener('click', toggleBooks);
recBookDiv.addEventListener('click', toggleBooks);



document.addEventListener('DOMContentLoaded', () => {
    const authButton = document.getElementById('authButton');
    
    // 检查登录状态
    if (localStorage.getItem('isLoggedIn') === 'true') {
        authButton.textContent = '退出登录';
        authButton.href = '#'; // 将 href 设置为 '#'，避免默认跳转
        
        authButton.addEventListener('click', (event) => {
            event.preventDefault(); // 阻止默认行为
            logout(); // 调用 logout 函数
            // 调用登出接口
            fetch('/ulogout', { method: 'POST' })
                .then(response => {
                    if (response.ok) {
                        // 清除登录状态
                        localStorage.removeItem('isLoggedIn');
                        window.location.href = '/index';
                    } else {
                        alert('登出失败，请重试。');
                    }
                })
                .catch(error => {
                    console.error('登出请求失败:', error);
                    alert('登出请求失败，请检查网络连接。');
                });
        });
    } else {
        authButton.textContent = '登录';
        authButton.href = '/login';
    }
});

// logout.js
function logout() {
    // 清除登录状态
    localStorage.removeItem('isLoggedIn');
}




