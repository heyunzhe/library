document.addEventListener('DOMContentLoaded', function() {
    const editButton = document.getElementById('editButton');
    const modal = document.getElementById('editModal');
    const closeModal = document.getElementById('closeModal');
    const editForm = document.getElementById('editForm');
    const changePasswordBtn = document.getElementById('changePasswordBtn');
    const passwordFields = document.getElementById('passwordFields');
    const returnModal = document.getElementById('model2');
    const returnForm = document.getElementById('returnForm');
    const submitReturnBtn = document.getElementById('submitReturn');
    let currentBookToReturn = '';
    let currentBookISBN = '';

    function fetchUserData() {
        fetch('/user/library', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            var userInfo = data.user_info;

            document.getElementById('user_name').textContent = userInfo.name;
            document.querySelector('.user-info p:nth-child(2)').textContent = `学生卡号: ${userInfo.username}`;
            document.getElementById('edit_name').value = userInfo.name;
            document.getElementById('edit_age').value = userInfo.age;
            document.getElementById('edit_birthday').value = userInfo.birthday;
            
            const userAvatar = document.querySelector('.user-avatar');
            if (userInfo.photo) {
                userAvatar.src = '../' + userInfo.photo;
            } else {
                userAvatar.src = '../userphoto/default-avatar.jpg';
            }

            const statCards = document.querySelectorAll('.stat-card');
            statCards[0].querySelector('p:nth-child(2)').textContent = `当前借阅: ${userInfo.user_cur_lend_amount}本`;
            statCards[0].querySelector('p:nth-child(3)').textContent = `历史借阅: ${userInfo.user_his_lend_amount}本`;
            statCards[1].querySelector('p').textContent = `已借阅: ${userInfo.user_cur_lend_amount}/5本`;

            
            const userInfoDiv = document.querySelector('.user-info div');
            const birthdayPara = document.createElement('p');
            birthdayPara.textContent = `生日: ${userInfo.birthday}`;
            const agePara = document.createElement('p');
            agePara.textContent = `年龄: ${userInfo.age}`;
            userInfoDiv.appendChild(birthdayPara);
            userInfoDiv.appendChild(agePara);


            
            const borrowHistoryList = document.getElementById('borrow_history_list');
            borrowHistoryList.innerHTML = ''; // 清空现有内容
            if (data.loan_history && data.loan_history.length > 0) {

                data.loan_history.forEach(Loan => {
                    const rows = document.createElement('tr');
                    rows.innerHTML = `
                    <td>${Loan.title}</td>
                    <td>${Loan.isbn}</td>
                    <td>${Loan.lend_date}</td>
                    <td>${Loan.return_date}</td>
                    <td>${Loan.late_fee}</td>
                    `;
                    borrowHistoryList.appendChild(rows);
                });
            }else{
                 borrowHistoryList.innerHTML = '<tr><td colspan="5">暂无借阅历史</td></tr>';
            }
            
            
            const currentBorrowList = document.getElementById('current_borrow_list');
            currentBorrowList.innerHTML = ''; // 清空现有内容
            if (data.current_loans && data.current_loans.length > 0) {
                data.current_loans.forEach(loan => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                    <td>${loan.title}</td>
                    <td>${loan.isbn}</td>
                    <td>${loan.lend_date}</td>
                    <td>${loan.exp_return_date}</td>
                        <td><button class="return-btn" data-book-name="${loan.title}" data-book-isbn="${loan.isbn}">归还</button></td>
                    `;
                    currentBorrowList.appendChild(row);
                });
            }else{
                 currentBorrowList.innerHTML = '<tr><td colspan="5">当前没有借阅的书籍</td></tr>';
            }
            
           
            // 移除硬编码的邮箱
            // const emailPara = userInfoDiv.querySelector('p:last-child');
            // if (emailPara && emailPara.textContent.includes('@')) {
            //     userInfoDiv.removeChild(emailPara);
            // }


            // 为归还按钮添加事件监听器
            document.querySelectorAll('.return-btn').forEach(btn => {
                btn.addEventListener('click', function() {
                    currentBookToReturn = this.getAttribute('data-book-name');
                    currentBookISBN = this.getAttribute('data-book-isbn');
                    document.getElementById('bookISBN').value = currentBookISBN;
                    returnModal.style.display = 'block';
                });
            });
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
            // alert('获取用户数据失败，请刷新页面重试。');
        });
    }

    fetchUserData();

    editButton.addEventListener('click', function() {
        modal.style.display = 'block';
    });

    closeModal.addEventListener('click', function() {
        modal.style.display = 'none';
    });

    changePasswordBtn.addEventListener('click', function() {
        passwordFields.style.display = passwordFields.style.display === 'none' ? 'block' : 'none';
    });

    editForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const formData = new FormData(editForm);


        fetch('/update/user', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.text();
        })
        .then(() => {
            alert('个人信息更新成功');
            modal.style.display = 'none';
            window.location.href = '/user/library';
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
            alert('更新失败，请重试');
        });
    });

    // 归还书籍模态窗口相关功能
    returnModal.querySelector('.close').addEventListener('click', function() {
        returnModal.style.display = 'none';
    });

    returnModal.querySelector('.btn-secondary').addEventListener('click', function() {
        returnModal.style.display = 'none';
    });

    submitReturnBtn.addEventListener('click', function(e) {
        e.preventDefault();
        const formData = new FormData(returnForm);
        // 不需要手动添加ISBN，因为它已经在隐藏字段中了

        fetch('/return/book', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (response.ok) {
                alert('归还成功');
                returnModal.style.display = 'none';
                window.location.href = '/user/library'; 
            } else if (response.status === 403) {
                alert('输入卡密错误');
            } else {
                throw new Error('Server error');
            }
        })
        .catch(error => {
            console.error('There was a problem with the return operation:', error);
            alert('服务器错误，请重试');
        });
    });
});