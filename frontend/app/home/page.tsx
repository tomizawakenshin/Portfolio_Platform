'use client';

import React, { useEffect, useState } from 'react';
import { User } from '../types/User';
import MnimumUserInfoInputModal from '../components/MinimumUserInfoInput';
import { useRouter } from 'next/navigation';

const HomePage = () => {
    const router = useRouter();
    const [user, setUser] = useState<User | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);

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
                if (!data.user.FirstName) {
                    setIsModalOpen(true); // FirstNameがnullならモーダルを開く
                }
            })
            .catch(error => {
                console.error('Error fetching user info:', error);
            });
    }, []);

    const handleModalClose = () => {
        setIsModalOpen(false);
    };

    const handleUserInfoSubmit = (firstName: string, lastName: string) => {
        // ここでAPIを呼び出してFirstNameとLastNameを設定します。
        fetch('http://localhost:8080/user/UpdateMinimumUserInfo', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include', // クッキーを含める
            body: JSON.stringify({ firstName, lastName }),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                // ユーザー情報を更新
                setUser({ ...user, FirstName: firstName, LastName: lastName } as User);
                setIsModalOpen(false);
            })
            .catch(error => {
                console.error('Error updating user info:', error);
            });
    };

    const handleLogout = () => {
        fetch('http://localhost:8080/auth/logout', {
            method: 'POST',
            credentials: 'include',
        })
            .then(response => {
                if (response.ok) {
                    // ログアウト成功時の処理
                    router.push("/auth"); // ログインページにリダイレクト
                } else {
                    throw new Error('Failed to logout');
                }
            })
            .catch(error => {
                console.error('Error during logout:', error);
            });
    };

    if (!user) {
        return <div>読み込み中...</div>;
    }

    return (
        <div>
            {/* ユーザーの情報があれば表示 */}
            <p>Welcome, {user.FirstName ? user.FirstName : 'Guest'}</p>

            {/* ログアウトボタンを追加 */}
            <button onClick={handleLogout}>ログアウト</button>

            {/* モーダルを表示 */}
            <MnimumUserInfoInputModal
                isOpen={isModalOpen}
                onClose={handleModalClose}
                onSubmit={handleUserInfoSubmit}
            />
        </div>
    );
};

export default HomePage;
