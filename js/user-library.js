document.addEventListener('DOMContentLoaded', function() {
    const $ = id => document.getElementById(id);
    const editModal = $('editModal');
    const returnModal = $('returnModal');
    const editForm = $('editForm');
    const returnForm = $('returnForm');

    // ==== 头像选择器 ====
    let selectedAvatarPath = '';

    document.querySelectorAll('.avatar-opt').forEach(opt => {
        opt.addEventListener('click', function() {
            const path = this.dataset.avatar;
            if (path === '') {
                // 上传按钮 → 打开文件选择
                $('editPhoto').click();
                return;
            }
            document.querySelectorAll('.avatar-opt').forEach(el => el.classList.remove('selected'));
            this.classList.add('selected');
            selectedAvatarPath = path;
            $('editAvatarPreview').src = path;
        });
    });

    $('editPhoto').addEventListener('change', function() {
        if (this.files && this.files[0]) {
            const reader = new FileReader();
            reader.onload = e => {
                $('editAvatarPreview').src = e.target.result;
                document.querySelectorAll('.avatar-opt').forEach(el => el.classList.remove('selected'));
                selectedAvatarPath = ''; // 使用上传的文件
            };
            reader.readAsDataURL(this.files[0]);
        }
    });

    // ==== 加载用户数据 ====
    fetch('/user/library', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
    })
    .then(r => { if (!r.ok) throw new Error('请求失败'); return r.json(); })
    .then(data => {
        const u = data.user_info;
        $('userName').textContent = u.name || u.username;
        $('userEmail').textContent = u.email || '--';
        $('userId').textContent = u.username;
        $('userAvatar').src = u.photo ? '/' + u.photo : '/images/default-avatar.svg';
        $('statCur').textContent = u.user_cur_lend_amount;
        $('statHis').textContent = u.user_his_lend_amount;
        $('statLimit').textContent = u.user_cur_lend_amount + '/5';
        $('statAge').textContent = u.age || '--';

        $('editName').value = u.name || '';
        $('editEmail').value = u.email || '';
        $('editBirthday').value = u.birthday || '';
        $('editAge').value = u.age || 0;
        $('editAvatarPreview').src = u.photo ? '/' + u.photo : '/images/default-avatar.svg';

        // 借阅历史
        const historyTbody = $('borrowHistoryList');
        historyTbody.innerHTML = '';
        if (data.loan_history && data.loan_history.length) {
            $('hisCount').textContent = data.loan_history.length + ' 条';
            data.loan_history.forEach(h => {
                const tr = document.createElement('tr');
                tr.innerHTML = `<td>${h.title}</td><td>${h.isbn}</td><td>${h.lend_date}</td><td>${h.return_date}</td><td>${h.late_fee || '0'}</td>`;
                historyTbody.appendChild(tr);
            });
        } else {
            historyTbody.innerHTML = '<tr><td colspan="5" style="text-align:center;color:#94a3b8;padding:24px">暂无借阅历史</td></tr>';
        }

        // 当前借阅
        const curTbody = $('currentBorrowList');
        curTbody.innerHTML = '';
        if (data.current_loans && data.current_loans.length) {
            $('curCount').textContent = data.current_loans.length + ' 本';
            data.current_loans.forEach(loan => {
                const tr = document.createElement('tr');
                tr.innerHTML = `<td>${loan.title}</td><td>${loan.isbn}</td><td>${loan.lend_date}</td><td>${loan.exp_return_date}</td>
                    <td style="display:flex;gap:6px">
                        <button class="read-btn" data-isbn="${loan.isbn}">阅读</button>
                        <button class="return-btn" data-isbn="${loan.isbn}" data-title="${loan.title}">归还</button>
                    </td>`;
                curTbody.appendChild(tr);
            });
            curTbody.querySelectorAll('.return-btn').forEach(btn => {
                btn.addEventListener('click', function() {
                    $('returnBookName').textContent = this.dataset.title;
                    $('bookISBN').value = this.dataset.isbn;
                    returnModal.classList.remove('hidden');
                });
            });
            curTbody.querySelectorAll('.read-btn').forEach(btn => {
                btn.addEventListener('click', function() {
                    window.location.href = '/read/book?isbn=' + this.dataset.isbn;
                });
            });
        } else {
            curTbody.innerHTML = '<tr><td colspan="5" style="text-align:center;color:#94a3b8;padding:24px">当前没有借阅的书籍</td></tr>';
        }
    })
    .catch(err => console.error('加载用户数据失败:', err));

    // ==== 生日自动算年龄 ====
    $('editBirthday').addEventListener('change', function() {
        if (!this.value) return;
        const b = new Date(this.value);
        const now = new Date();
        let age = now.getFullYear() - b.getFullYear();
        const mDiff = now.getMonth() - b.getMonth();
        if (mDiff < 0 || (mDiff === 0 && now.getDate() < b.getDate())) age--;
        $('editAge').value = age >= 0 ? age : 0;
    });

    // ==== 弹窗控制 ====
    $('editButton').addEventListener('click', () => editModal.classList.remove('hidden'));
    $('closeModal').addEventListener('click', () => editModal.classList.add('hidden'));
    $('cancelEdit').addEventListener('click', () => editModal.classList.add('hidden'));
    $('closeReturn').addEventListener('click', () => returnModal.classList.add('hidden'));
    $('cancelReturn').addEventListener('click', () => returnModal.classList.add('hidden'));

    editModal.addEventListener('click', e => { if (e.target === editModal) editModal.classList.add('hidden'); });
    returnModal.addEventListener('click', e => { if (e.target === returnModal) returnModal.classList.add('hidden'); });

    // ==== 提交编辑 ====
    editForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const fd = new FormData();
        fd.append('name', $('editName').value);
        fd.append('birthday', $('editBirthday').value);

        const pwd = $('editPassword').value;
        const newPwd = $('editNewPassword').value;
        if (pwd) fd.append('password', pwd);
        if (newPwd) fd.append('newpassword', newPwd);

        // 头像：优先用上传的文件，否则用选择的默认头像
        const photoFile = $('editPhoto').files[0];
        if (photoFile) {
            fd.append('photo', photoFile);
        } else if (selectedAvatarPath) {
            // 发送选择的默认头像路径
            fd.append('avatar_path', selectedAvatarPath);
        }

        fetch('/update/user', { method: 'POST', body: fd })
            .then(r => r.json().catch(() => r.text()))
            .then(res => {
                if (res && res.error) { alert(res.error); return; }
                alert('保存成功');
                editModal.classList.add('hidden');
                window.location.reload();
            })
            .catch(err => { alert('更新失败'); console.error(err); });
    });

    // ==== 归还书籍 ====
    returnForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const fd = new FormData(returnForm);
        fetch('/return/book', { method: 'POST', body: fd })
            .then(r => {
                if (r.ok) { alert('归还成功'); returnModal.classList.add('hidden'); window.location.reload(); }
                else if (r.status === 403) { alert('卡密错误'); }
                else throw new Error('服务器错误');
            })
            .catch(err => { alert('归还失败'); console.error(err); });
    });

    // ==== 退出登录 ====
    $('logoutLink').addEventListener('click', function(e) {
        e.preventDefault();
        fetch('/ulogout', { method: 'POST' })
            .then(() => { localStorage.removeItem('isLoggedIn'); window.location.href = '/index'; })
            .catch(() => { localStorage.removeItem('isLoggedIn'); window.location.href = '/index'; });
    });
});
