document.addEventListener('DOMContentLoaded', function() {
    const rankingList = document.getElementById('ranking-list');

    fetch('/ranking', {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        // 按借阅次数降序排序
        data.sort((a, b) => b.count - a.count);

        data.forEach((book, index) => {
            const rankingItem = document.createElement('div');
            rankingItem.className = 'ranking-item';
            
            const rankClass = index < 3 ? `rank-${index + 1}` : '';
            
            rankingItem.innerHTML = `
                <div class="rank ${rankClass}">${index + 1}</div>
                <img src="../images/image${book.isbn}.jpg" alt="${book.title}" class="book-cover">
                <div class="book-info">
                    <div class="book-title">${book.title}</div>
                    <div class="book-isbn">ISBN: ${book.isbn}</div>
                    <div class="borrow-count">借阅次数: ${book.count}</div>
                </div>
            `;
            
            rankingList.appendChild(rankingItem);
        });
    })
    .catch(error => {
        console.error('Error:', error);
        rankingList.innerHTML = '<p>加载排行榜数据时出错，请稍后再试。</p>';
    });
});