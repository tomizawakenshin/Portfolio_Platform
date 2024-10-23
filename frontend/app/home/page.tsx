// home/page.tsx

'use client';

import React, { useEffect, useState } from 'react';
import { User } from '../types/User';
import MnimumUserInfoInputModal from '../components/MinimumUserInfoInput';
import { useRouter } from 'next/navigation';
import useAuthCheck from '../hooks/useAuthCheck';
import WelcomeModal from '../components/WelcomeModal';

const HomePage = () => {
    useAuthCheck(); // ログインチェック

    const router = useRouter();
    const [user, setUser] = useState<User | null>(null);
    const [isWelcomeModalOpen, setIsWelcomeModalOpen] = useState(false);
    const [isMinimumUserInputModalOpen, setIsMinimumUserInputModalOpen] = useState(false);

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
                    // FirstNameがない場合は WelcomeModal を表示
                    setIsWelcomeModalOpen(true);
                }
            })
            .catch(error => {
                console.error('Error fetching user info:', error);
            });
    }, []);

    const handleWelcomeModalStart = () => {
        setIsWelcomeModalOpen(false);
        setIsMinimumUserInputModalOpen(true);
    };

    const handleMinimumUserInputModalClose = () => {
        setIsMinimumUserInputModalOpen(false);
    };

    const handleUserInfoSubmit = (
        firstName: string,
        lastName: string,
        firstNameKana: string,
        lastNameKana: string,
        schoolName: string,
        department: string,
        laboratory: string,
        graduationYear: string,
        desiredJobTypes: string[],
        skills: string[] // スキルも追加
    ) => {
        // APIを呼び出してユーザー情報を更新
        fetch('http://localhost:8080/user/UpdateMinimumUserInfo', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include', // クッキーを含める
            body: JSON.stringify({
                firstName,
                lastName,
                firstNameKana,
                lastNameKana,
                schoolName,
                department,
                laboratory,
                graduationYear,
                desiredJobTypes,
                skills,
            }),
        })
            .then(response => {
                if (!response.ok) {
                    return response.json().then((data) => {
                        throw new Error(data.error || 'Network response was not ok');
                    });
                }
                return response.json();
            })
            .then(data => {
                // ユーザー情報を更新
                setUser({
                    ...user,
                    FirstName: firstName,
                    LastName: lastName,
                    FirstNameKana: firstNameKana,
                    LastNameKana: lastNameKana,
                    SchoolName: schoolName,
                    Department: department,
                    Laboratory: laboratory,
                    GraduationYear: graduationYear,
                    DesiredJobTypes: desiredJobTypes,
                    Skills: skills,
                } as User);
                setIsMinimumUserInputModalOpen(false);
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

            {/* WelcomeModalを表示 */}
            <WelcomeModal
                isOpen={isWelcomeModalOpen}
                onClose={() => setIsWelcomeModalOpen(false)}
                onStart={handleWelcomeModalStart}
            />

            {/* MinimumUserInputModalを表示 */}
            <MnimumUserInfoInputModal
                isOpen={isMinimumUserInputModalOpen}
                onClose={handleMinimumUserInputModalClose}
                onSubmit={handleUserInfoSubmit}
            />
        </div>
    );
};

export default HomePage;
