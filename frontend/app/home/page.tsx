"use client";

import React, { useEffect, useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { User } from '../types/User';
import { Portfolio } from '../types/Portfolio';
import useAuthCheck from '../hooks/useAuthCheck';

import Header from '../components/Header_Home'; // ← 作成したHeaderコンポーネントをimport
import WelcomeModal from '../components/WelcomeModal';
import MnimumUserInfoInputModal from '../components/MinimumUserInfoInput';

// カタカナ→イニシャル（英字）の簡易変換例 (先頭1文字だけのマッピング)
function getInitial(kanaChar: string): string {
    // 先頭一文字を取り出す
    const firstChar = kanaChar.slice(0, 1);
    // 簡易辞書
    const dictionary: Record<string, string> = {
        'ア': 'A', 'イ': 'I', 'ウ': 'U', 'エ': 'E', 'オ': 'O',
        'カ': 'K', 'キ': 'K', 'ク': 'K', 'ケ': 'K', 'コ': 'K',
        'サ': 'S', 'シ': 'S', 'ス': 'S', 'セ': 'S', 'ソ': 'S',
        'タ': 'T', 'チ': 'T', 'ツ': 'T', 'テ': 'T', 'ト': 'T',
        'ナ': 'N', 'ニ': 'N', 'ヌ': 'N', 'ネ': 'N', 'ノ': 'N',
        'ハ': 'H', 'ヒ': 'H', 'フ': 'H', 'ヘ': 'H', 'ホ': 'H',
        'マ': 'M', 'ミ': 'M', 'ム': 'M', 'メ': 'M', 'モ': 'M',
        'ヤ': 'Y', 'ユ': 'Y', 'ヨ': 'Y',
        'ラ': 'R', 'リ': 'R', 'ル': 'R', 'レ': 'R', 'ロ': 'R',
        'ワ': 'W', 'ヲ': 'W', 'ン': 'N',
    };
    return dictionary[firstChar] || 'X'; // 該当がなければ X とか適当に
}

// ユーザーのカタカナからイニシャルを取得
function getInitialsFromKana(firstNameKana: string, lastNameKana: string): string {
    // 例: LastNameKana="タニ", FirstNameKana="ケンシン" => "T.K."
    const lastInit = getInitial(lastNameKana);
    const firstInit = getInitial(firstNameKana);
    return `${lastInit}.${firstInit}.`;
}

const HomePage = () => {
    useAuthCheck(); // ログインチェック
    const router = useRouter();

    // -----------------------------
    // ユーザー・作品関連
    // -----------------------------
    const [user, setUser] = useState<User | null>(null);
    const [portfolio, setPortfolio] = useState<Portfolio[]>([]);

    // -----------------------------
    // モーダル関係
    // -----------------------------
    const [isWelcomeModalOpen, setIsWelcomeModalOpen] = useState(false);
    const [isMinimumUserInputModalOpen, setIsMinimumUserInputModalOpen] = useState(false);

    // -----------------------------
    // 検索バーUIの状態
    // -----------------------------
    const [searchKeyword, setSearchKeyword] = useState("");

    // 卒業年 (ドロップダウン)
    const [isYearDropdownOpen, setIsYearDropdownOpen] = useState(false);
    const yearDropdownRef = useRef<HTMLDivElement>(null);
    const [selectedYears, setSelectedYears] = useState<string[]>([]);
    const [graduationYearOptions, setGraduationYearOptions] = useState<string[]>([]);

    // ジャンル (ドロップダウン)
    const [isGenreDropdownOpen, setIsGenreDropdownOpen] = useState(false);
    const genreDropdownRef = useRef<HTMLDivElement>(null);
    const [availableGenres, setAvailableGenres] = useState<string[]>([]);
    const [selectedGenres, setSelectedGenres] = useState<string[]>([]);

    // スキル (モーダル)
    const [isSkillModalOpen, setIsSkillModalOpen] = useState(false);
    const [availableSkills, setAvailableSkills] = useState<string[]>([]);
    const [skillInput, setSkillInput] = useState("");
    const [suggestedSkills, setSuggestedSkills] = useState<string[]>([]);
    const [selectedSkills, setSelectedSkills] = useState<string[]>([]);
    const skillContainerRef = useRef<HTMLDivElement>(null);

    // =========================
    // 卒業年を動的に取得する関数
    // =========================
    const getGraduationYearOptions = () => {
        const options: string[] = [];
        const today = new Date();
        const currentYear = today.getFullYear();
        const month = today.getMonth();
        const day = today.getDate();

        // 3月31日以前なら今年含め4年分、それ以降なら来年始まりの4年分
        const isBeforeMarch31 = month < 2 || (month === 2 && day <= 30);
        if (isBeforeMarch31) {
            for (let i = 0; i < 4; i++) {
                options.push(`${currentYear + i}卒`);
            }
        } else {
            for (let i = 1; i <= 4; i++) {
                options.push(`${currentYear + i}卒`);
            }
        }
        return options;
    };

    // -----------------------------
    // 初期データの取得
    // -----------------------------
    useEffect(() => {
        // 卒業年のオプションをセット
        setGraduationYearOptions(getGraduationYearOptions());

        // ユーザー情報
        fetch('http://localhost:8080/user/GetInfo', { credentials: 'include' })
            .then(res => {
                if (!res.ok) throw new Error("Failed to fetch user info");
                return res.json();
            })
            .then(data => {
                setUser(data.user);
                if (!data.user.FirstName) {
                    setIsWelcomeModalOpen(true);
                }
            })
            .catch(err => console.error(err));

        // 作品情報
        fetch('http://localhost:8080/Portfolio/getAllPosts', { credentials: 'include' })
            .then(res => res.json())
            .then(data => {
                setPortfolio(data.portfolio);
            })
            .catch(err => console.error(err));

        // ジャンル一覧
        fetch('http://localhost:8080/options/genre', { credentials: 'include' })
            .then(res => res.json())
            .then(data => {
                setAvailableGenres(data.genres);
            })
            .catch(err => console.error(err));

        // スキル一覧
        fetch('http://localhost:8080/options/skills', { credentials: 'include' })
            .then(res => res.json())
            .then(data => {
                setAvailableSkills(data.skills);
            })
            .catch(err => console.error(err));
    }, []);

    // -----------------------------
    // ドロップダウンの外側クリック
    // -----------------------------
    // 卒業年
    useEffect(() => {
        function handleClickOutside(e: MouseEvent) {
            if (
                yearDropdownRef.current &&
                !yearDropdownRef.current.contains(e.target as Node)
            ) {
                setIsYearDropdownOpen(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    // ジャンル
    useEffect(() => {
        function handleClickOutside(e: MouseEvent) {
            if (
                genreDropdownRef.current &&
                !genreDropdownRef.current.contains(e.target as Node)
            ) {
                setIsGenreDropdownOpen(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    // -----------------------------
    // スキルモーダル: 候補絞り込み
    // -----------------------------
    useEffect(() => {
        if (skillInput) {
            const filtered = availableSkills.filter(
                (skill) =>
                    skill.toLowerCase().includes(skillInput.toLowerCase()) &&
                    !selectedSkills.includes(skill)
            );
            setSuggestedSkills(filtered);
        } else {
            setSuggestedSkills([]);
        }
    }, [skillInput, selectedSkills, availableSkills]);

    // -----------------------------
    // スキルモーダル: 外側クリック
    // -----------------------------
    useEffect(() => {
        function handleSkillModalClickOutside(e: MouseEvent) {
            if (
                skillContainerRef.current &&
                !skillContainerRef.current.contains(e.target as Node)
            ) {
                // モーダル外クリック時の挙動は必要に応じて。
            }
        }
        document.addEventListener("mousedown", handleSkillModalClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleSkillModalClickOutside);
        };
    }, []);

    // -----------------------------
    // イベントハンドラ
    // -----------------------------
    const handlePost = () => {
        router.push("/post");
    };

    const handleLogout = () => {
        fetch('http://localhost:8080/auth/logout', { method: "POST", credentials: 'include' })
            .then(res => {
                if (res.ok) {
                    router.push("/auth");
                } else {
                    throw new Error("Failed to logout");
                }
            })
            .catch(err => console.error(err));
    };

    // WelcomeModal -> MinimumUserInfoInputModal
    const handleWelcomeModalStart = () => {
        setIsWelcomeModalOpen(false);
        setIsMinimumUserInputModalOpen(true);
    };
    const handleMinimumUserInputModalClose = () => {
        setIsMinimumUserInputModalOpen(false);
    };

    // ユーザー情報更新
    const handleUserInfoSubmit = async (
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
        try {
            // PUT リクエストでサーバーに更新を依頼
            const res = await fetch('http://localhost:8080/user/UpdateMinimumUserInfo', {
                method: 'PUT',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    FirstName: firstName,
                    LastName: lastName,
                    FirstNameKana: firstNameKana,
                    LastNameKana: lastNameKana,
                    SchoolName: schoolName,
                    Department: department,
                    Laboratory: laboratory,
                    // 「2025卒」のような文字列の場合、もしサーバーが数値を期待するなら変換処理が必要
                    GraduationYear: graduationYear.replace("卒", ""),
                    DesiredJobTypes: desiredJobTypes,
                    Skills: skills,
                }),
            });

            if (!res.ok) {
                throw new Error('ユーザー情報の更新に失敗しました');
            }

            const data = await res.json();
            // サーバーが更新後のユーザー情報を返す場合
            if (data.user) {
                // 更新したユーザーをステートに反映
                setUser(data.user);
            }

            // モーダルを閉じる
            setIsMinimumUserInputModalOpen(false);

        } catch (error) {
            console.error(error);
            alert("ユーザー情報の更新中にエラーが発生しました。");
        }
    };

    // 卒業年ドロップダウン
    const toggleYearDropdown = () => {
        setIsYearDropdownOpen(!isYearDropdownOpen);
    };
    const handleYearCheck = (year: string) => {
        if (selectedYears.includes(year)) {
            setSelectedYears(selectedYears.filter((y) => y !== year));
        } else {
            setSelectedYears([...selectedYears, year]);
        }
    };

    // ジャンルドロップダウン
    const toggleGenreDropdown = () => {
        setIsGenreDropdownOpen(!isGenreDropdownOpen);
    };
    const handleGenreCheck = (genre: string) => {
        if (selectedGenres.includes(genre)) {
            setSelectedGenres(selectedGenres.filter((g) => g !== genre));
        } else {
            setSelectedGenres([...selectedGenres, genre]);
        }
    };

    // スキルモーダル
    const openSkillModal = () => {
        setIsSkillModalOpen(true);
    };
    const closeSkillModal = () => {
        setIsSkillModalOpen(false);
        setSkillInput("");
        setSuggestedSkills([]);
    };
    const handleSkillSelect = (skill: string) => {
        if (selectedSkills.length >= 5) {
            alert("スキルは最大5件まで選択できます。");
            return;
        }
        if (!selectedSkills.includes(skill)) {
            setSelectedSkills([...selectedSkills, skill]);
        }
        setSkillInput("");
        setSuggestedSkills([]);
    };
    const removeSkill = (skill: string) => {
        setSelectedSkills(selectedSkills.filter((s) => s !== skill));
    };
    const handleSkillInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            e.preventDefault();
            if (skillInput.trim()) {
                handleSkillSelect(skillInput.trim());
            }
        }
    };

    // 作品カードをクリック -> 詳細ページへ
    const handlePortfolioClick = (postId: number) => {
        router.push(`/Portfolio/${postId}`);
    };

    // -----------------------------
    // フィルタロジック
    // -----------------------------
    const filteredPortfolio = portfolio.filter((post) => {
        // キーワード
        const keyword = searchKeyword.toLowerCase();
        const matchesKeyword =
            !keyword ||
            post.Title.toLowerCase().includes(keyword) ||
            post.Description.toLowerCase().includes(keyword);

        // 卒業年
        let matchesYear = true;
        if (selectedYears.length > 0) {
            if (!post.User?.GraduationYear) {
                matchesYear = false;
            } else {
                const userYear = post.User.GraduationYear.toString();
                matchesYear = selectedYears.includes(userYear);
            }
        }

        // ジャンル (OR検索 => some)
        let matchesGenre =
            selectedGenres.length === 0 ||
            (post.Genres && selectedGenres.some((g) => post.Genres.includes(g)));

        // スキル (OR検索 => some)
        let matchesSkill =
            selectedSkills.length === 0 ||
            (post.Skills && selectedSkills.some((s) => post.Skills.includes(s)));

        return (
            matchesKeyword &&
            matchesYear &&
            matchesGenre &&
            matchesSkill
        );
    });

    // -----------------------------
    // レンダリング
    // -----------------------------
    if (!user) {
        return <div>読み込み中...</div>;
    }

    // 写真の有無をHeaderに渡す
    const userHasPhoto = !!user.ProfileImageURL;
    const userPhotoURL = user.ProfileImageURL
        ? `http://localhost:8080/${user.ProfileImageURL}`
        : undefined;

    return (
        <div className="min-h-screen bg-gray-50">
            {/* ヘッダー */}
            <Header
                onPostClick={handlePost}
                userHasPhoto={userHasPhoto}
                userPhotoURL={userPhotoURL}
            />

            <div className="p-4">
                {/* 検索バー全体をカード風に */}
                <div className="bg-white shadow-md rounded p-4 mb-6">
                    <div className="flex items-center space-x-2">
                        {/* キーワード検索を短めに */}
                        <div className="w-64">
                            <input
                                type="text"
                                placeholder="キーワードで検索"
                                className="w-full border px-2 py-1 rounded"
                                value={searchKeyword}
                                onChange={(e) => setSearchKeyword(e.target.value)}
                            />
                        </div>

                        {/* 卒業年 (ドロップダウン) */}
                        <div className="relative" ref={yearDropdownRef}>
                            <button
                                onClick={toggleYearDropdown}
                                className={`px-3 py-2 border rounded ${selectedYears.length > 0
                                    ? "bg-orange-500 text-white"
                                    : "bg-white text-black"
                                    }`}
                            >
                                卒業年 ▼
                            </button>
                            {isYearDropdownOpen && (
                                <div className="absolute left-0 mt-2 bg-white border rounded p-2 z-50">
                                    {graduationYearOptions.map((year) => (
                                        <label
                                            key={year}
                                            className="flex items-center cursor-pointer mb-1"
                                        >
                                            <input
                                                type="checkbox"
                                                className="mr-2"
                                                checked={selectedYears.includes(year)}
                                                onChange={() => handleYearCheck(year)}
                                            />
                                            {year}
                                        </label>
                                    ))}
                                </div>
                            )}
                        </div>

                        {/* ジャンル (ドロップダウン) */}
                        <div className="relative" ref={genreDropdownRef}>
                            <button
                                onClick={toggleGenreDropdown}
                                className={`px-3 py-2 border rounded ${selectedGenres.length > 0
                                    ? "bg-orange-500 text-white"
                                    : "bg-white text-black"
                                    }`}
                            >
                                ジャンル ▼
                            </button>
                            {isGenreDropdownOpen && (
                                <div className="absolute left-0 mt-2 bg-white border rounded p-2 z-50 w-48 max-h-60 overflow-y-auto">
                                    {availableGenres.map((genre) => (
                                        <label
                                            key={genre}
                                            className="flex items-center cursor-pointer mb-1"
                                        >
                                            <input
                                                type="checkbox"
                                                className="mr-2"
                                                checked={selectedGenres.includes(genre)}
                                                onChange={() => handleGenreCheck(genre)}
                                            />
                                            {genre}
                                        </label>
                                    ))}
                                </div>
                            )}
                        </div>

                        {/* スキル (モーダルオープン) */}
                        <div>
                            <button
                                onClick={openSkillModal}
                                className={`px-3 py-2 border rounded ${selectedSkills.length > 0
                                    ? "bg-orange-500 text-white"
                                    : "bg-white text-black"
                                    }`}
                            >
                                スキル ▼
                            </button>
                        </div>
                    </div>
                </div>

                {/* スキルモーダル */}
                {isSkillModalOpen && (
                    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                        <div
                            className="bg-white rounded p-6 relative max-w-md w-full"
                            ref={skillContainerRef}
                        >
                            <button
                                className="absolute top-2 right-2 text-gray-600"
                                onClick={closeSkillModal}
                            >
                                ×
                            </button>
                            <h2 className="text-xl font-semibold mb-2">
                                スキル (5件まで)
                            </h2>
                            {/* 選択済みのスキル一覧 */}
                            <div className="flex flex-wrap gap-2 mb-4">
                                {selectedSkills.map((skill) => (
                                    <span
                                        key={skill}
                                        className="bg-gray-200 px-3 py-1 rounded-full text-sm flex items-center"
                                    >
                                        {skill}
                                        <button
                                            className="ml-2 text-red-500"
                                            onClick={() => removeSkill(skill)}
                                        >
                                            &times;
                                        </button>
                                    </span>
                                ))}
                            </div>
                            {/* 検索用のテキストボックス */}
                            <input
                                type="text"
                                placeholder="スキルで検索"
                                value={skillInput}
                                onChange={(e) => setSkillInput(e.target.value)}
                                onKeyDown={handleSkillInputKeyDown}
                                className="w-full border px-3 py-2 rounded mb-2"
                            />
                            {/* 候補リスト */}
                            {skillInput && suggestedSkills.length > 0 && (
                                <div className="border rounded p-2 max-h-40 overflow-y-auto">
                                    {suggestedSkills.map((skill) => (
                                        <div
                                            key={skill}
                                            className="p-1 hover:bg-gray-100 cursor-pointer"
                                            onClick={() => handleSkillSelect(skill)}
                                        >
                                            {skill}
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {/* Welcome & MinimumUserInfo モーダル */}
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

                {/* 作品一覧 (フィルタ後) */}
                <h2 className="text-2xl font-bold mb-4">作品</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                    {filteredPortfolio.map((post) => {
                        // ユーザー写真 or デフォルト
                        const userPhoto = post.User?.ProfileImageURL
                            ? `http://localhost:8080/${post.User.ProfileImageURL}`
                            : "/images/defaultUserIcon.png";

                        // イニシャルを取得
                        let userInitials = "N/A";
                        if (post.User?.LastNameKana && post.User?.FirstNameKana) {
                            userInitials = getInitialsFromKana(
                                post.User.FirstNameKana,
                                post.User.LastNameKana
                            );
                        }

                        return (
                            <div
                                key={post.ID}
                                // カード全体にホバー時の影などをつける例
                                className="
                                bg-white shadow-md rounded p-4 cursor-pointer 
                                transition-shadow duration-200 hover:shadow-lg
                              "
                                onClick={() => handlePortfolioClick(post.ID)}
                            >
                                {/* 画像を拡大するために overflow-hidden + group を指定 */}
                                <div className="relative overflow-hidden rounded group">
                                    {post.Images && post.Images.length > 0 && (
                                        <img
                                            src={`http://localhost:8080/${post.Images[0].URL}`}
                                            alt={post.Title}
                                            // group-hover:scale-105 を指定し、ホバー時に拡大
                                            className="w-full h-40 object-cover transition-transform duration-200 group-hover:scale-105"
                                        />
                                    )}
                                    {/* お気に入りボタン */}
                                    <button
                                        className="absolute top-2 right-2 bg-white p-1 rounded-full"
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            console.log("お気に入りボタンの処理");
                                        }}
                                    >
                                        ❤️
                                    </button>
                                </div>

                                <h3 className="text-lg font-semibold mt-2">{post.Title}</h3>

                                <div className="flex items-center mt-2">
                                    <img
                                        src={userPhoto}
                                        alt="user icon"
                                        className="w-8 h-8 rounded-full mr-2 object-cover"
                                    />
                                    <p className="text-sm font-medium text-gray-700">
                                        {userInitials}
                                    </p>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>
        </div>
    );
};

export default HomePage;
