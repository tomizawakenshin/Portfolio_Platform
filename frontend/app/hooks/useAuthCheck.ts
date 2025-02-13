// hooks/useAuthCheck.ts
'use client';

import { BACKEND_URL } from '@/config';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

const useAuthCheck = () => {
    const router = useRouter();

    useEffect(() => {
        const checkAuth = async () => {
            try {
                const response = await fetch(`${BACKEND_URL}/auth/check`, {
                    method: 'GET',
                    credentials: 'include', // クッキーを含める
                });

                if (response.status === 200) {
                    // 認証されている場合は何もしない
                    return;
                } else {
                    // 認証されていない場合は /auth へリダイレクト
                    router.push('/auth');
                }
            } catch (error) {
                console.error('Error checking authentication:', error);
                // エラーが発生した場合も /auth へリダイレクト
                router.push('/auth');
            }
        };
        checkAuth();
    }, [router]);
};

export default useAuthCheck;
