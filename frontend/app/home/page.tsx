'use client'

import React, { useEffect, useState } from 'react';
import { User } from '../types/User';

const HomePage = () => {
    const [user, setUser] = useState<User | null>(null);

    useEffect(() => {
        fetch('http://localhost:8080/user/GetInfo', {
            credentials: 'include', // クッキーを含める
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                setUser(data.user);
            })
            .catch(error => {
                console.error('Error fetching user info:', error);
            });
    }, []);

    if (!user) {
        return <div>読み込み中...</div>;
    }

    return (
        <div>
            本登録が完了しました！
            <br />
            あなたのメールアドレス: {user.Email}
            <br />
            {/* 他のユーザー情報を表示する場合は、以下に追加 */}
            {/* 例: ユーザー名: {user.username} */}
        </div>
    );
};

export default HomePage;
