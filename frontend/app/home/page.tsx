'use client'

import React, { useEffect, useState } from 'react';
import { User } from '../types/User';
import MnimumUserInfoInputModal from '../components/MinimumUserInfoInput';
import { useRouter } from 'next/navigation';
import useAuthCheck from '../hooks/useAuthCheck';
import WelcomeModal from '../components/WelcomeModal';
import { Portfolio } from '../types/Portfolio';

const HomePage = () => {
    useAuthCheck(); // ログインチェック

    const router = useRouter();
    const [user, setUser] = useState<User | null>(null);
    const [isWelcomeModalOpen, setIsWelcomeModalOpen] = useState(false);
    const [isMinimumUserInputModalOpen, setIsMinimumUserInputModalOpen] = useState(false);
    const [portfolio, setPortfolio] = useState<Portfolio[]>([]);

    useEffect(() => {
        fetch('http://localhost:8080/user/GetInfo', {
            credentials: 'include',
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
                    setIsWelcomeModalOpen(true);
                }
            })
            .catch(error => {
                console.error('Error fetching user info:', error);
            });
    }, []);

    useEffect(() => {
        fetch('http://localhost:8080/Portfolio/getAllPosts', {
            credentials: 'include',
        })
            .then(response => response.json())
            .then(data => {
                setPortfolio(data.portfolio);
            })
            .catch(error => {
                console.error('Error fetching portfolio:', error);
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
        skills: string[]
    ) => {
        fetch('http://localhost:8080/user/UpdateMinimumUserInfo', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
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
                    router.push("/auth");
                } else {
                    throw new Error('Failed to logout');
                }
            })
            .catch(error => {
                console.error('Error during logout:', error);
            });
    };

    const handlePost = () => {
        router.push('/post/');
    };

    if (!user) {
        return <div>読み込み中...</div>;
    }

    return (
        <div className="p-8">
            <p className="text-lg mb-4">Welcome, {user.FirstName ? user.FirstName : 'Guest'}</p>

            <button onClick={handlePost} className="mb-4 px-4 py-2 bg-blue-500 text-white rounded">投稿する</button>
            <button onClick={handleLogout} className="ml-2 px-4 py-2 bg-red-500 text-white rounded">ログアウト</button>

            <WelcomeModal
                isOpen={isWelcomeModalOpen}
                onClose={() => setIsWelcomeModalOpen(false)}
                onStart={handleWelcomeModalStart}
            />

            <MnimumUserInfoInputModal
                isOpen={isMinimumUserInputModalOpen}
                onClose={handleMinimumUserInputModalClose}
                onSubmit={handleUserInfoSubmit}
            />

            <h2 className="text-2xl font-bold mt-8 mb-4">作品</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {portfolio.map((portfolio) => (
                    <div key={portfolio.ID} className="bg-white shadow-md rounded-lg p-4">
                        <div className="relative">
                            {portfolio.Images && portfolio.Images.length > 0 && (
                                <img
                                    src={`http://localhost:8080/${portfolio.Images[0].URL}`}
                                    alt={portfolio.Title}
                                    className="w-full h-40 object-cover rounded-md"
                                />
                            )}
                            <button className="absolute top-2 right-2 bg-white p-1 rounded-full">
                                <span role="img" aria-label="favorite">❤️</span>
                            </button>
                        </div>
                        <h3 className="text-lg font-semibold mt-2">{portfolio.Title}</h3>
                        <p className="text-gray-500 text-sm">{portfolio.Description}</p>
                        <div className="flex items-center mt-2">
                            {portfolio.User && (
                                <img
                                    src={`http://localhost:8080/${portfolio.User.profilePictureURL}`} // プロフィール画像がある場合
                                    alt={portfolio.User.FirstName}
                                    className="w-8 h-8 rounded-full mr-2"
                                />
                            )}
                            <p className="text-sm font-medium text-gray-700">{portfolio.User ? portfolio.User.FirstName : 'Unknown'}</p>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default HomePage;
