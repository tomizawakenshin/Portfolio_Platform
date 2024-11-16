import { useRouter } from 'next/navigation';
import React, { useEffect, useRef, useState } from 'react';

interface MnimumUserInfoInputModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: (
        firstName: string,
        lastName: string,
        firstNameKana: string,
        lastNameKana: string,
        schoolName: string,
        department: string,
        laboratory: string,
        graduationYear: string,
        desiredJobTypes: string[],
        skills: string[] // スキルを保持
    ) => void;
}

const MnimumUserInfoInputModal: React.FC<MnimumUserInfoInputModalProps> = ({ isOpen, onClose, onSubmit }) => {
    const router = useRouter();
    const dropdownRef = useRef<HTMLDivElement>(null);
    const skillContainerRef = useRef<HTMLDivElement>(null); // Ref for skills input container

    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [firstNameKana, setFirstNameKana] = useState('');
    const [lastNameKana, setLastNameKana] = useState('');
    const [schoolName, setSchoolName] = useState('');
    const [department, setDepartment] = useState('');
    const [laboratory, setLaboratory] = useState('');
    const [graduationYear, setGraduationYear] = useState('');
    const [desiredJobTypes, setDesiredJobTypes] = useState<string[]>([]); // 希望職種の状態を保持
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [skills, setSkills] = useState<string[]>([]); // Selected skills
    const [skillInput, setSkillInput] = useState(''); // For the input field
    const [suggestedSkills, setSuggestedSkills] = useState<string[]>([]);
    const [jobTypeOptions, setJobTypeOptions] = useState<string[]>([]);
    const [availableSkills, setAvailableSkills] = useState<string[]>([]);

    // 追加するステート
    const [errors, setErrors] = useState<{
        firstName?: string;
        lastName?: string;
        firstNameKana?: string;
        lastNameKana?: string;
        schoolName?: string;
        department?: string;
        laboratory?: string;
        graduationYear?: string;
        desiredJobTypes?: string;
    }>({});

    // カタカナの正規表現
    const katakanaRegex = /^[ァ-ンヴー]*$/;

    useEffect(() => {
        fetch('http://localhost:8080/options/job-types', {
            credentials: 'include',
        })
            .then((response) => response.json())
            .then((data) => {
                setJobTypeOptions(data.jobTypes);
            })
            .catch((error) => {
                console.error('Error fetching job types:', error);
            });

        fetch('http://localhost:8080/options/skills', {
            credentials: 'include',
        })
            .then((response) => response.json())
            .then((data) => {
                setAvailableSkills(data.skills);
            })
            .catch((error) => {
                console.error('Error fetching skills:', error);
            });
    }, []);

    useEffect(() => {
        // Close desired job types dropdown when clicking outside
        function handleClickOutside(event: MouseEvent) {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setIsDropdownOpen(false);
            }
        }
        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);

    useEffect(() => {
        // Close skills suggestion dropdown when clicking outside and clear input
        function handleSkillClickOutside(event: MouseEvent) {
            if (
                skillContainerRef.current &&
                !skillContainerRef.current.contains(event.target as Node)
            ) {
                setSuggestedSkills([]);
                setSkillInput('');
            }
        }
        document.addEventListener('mousedown', handleSkillClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleSkillClickOutside);
        };
    }, []);

    useEffect(() => {
        // Filter skills based on input
        if (skillInput) {
            const filtered = availableSkills.filter(skill =>
                skill.toLowerCase().includes(skillInput.toLowerCase()) && !skills.includes(skill)
            );
            setSuggestedSkills(filtered);
        } else {
            setSuggestedSkills([]);
        }
    }, [skillInput, skills, availableSkills]);

    const handleSkillSelect = (skill: string) => {
        if (!skills.includes(skill)) {
            setSkills([...skills, skill]);
        }
        setSkillInput('');
        setSuggestedSkills([]);
    };

    const removeSkill = (skillToRemove: string) => {
        setSkills(skills.filter(skill => skill !== skillToRemove));
    };

    const handleSkillInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            if (skillInput.trim()) {
                handleSkillSelect(skillInput.trim());
            }
        }
    };

    const getGraduationYearOptions = () => {
        const options = [];
        const today = new Date();
        const currentYear = today.getFullYear();
        const month = today.getMonth(); // 0-based (0 = January)
        const day = today.getDate();

        const isBeforeMarch31 = month < 2 || (month === 2 && day <= 30); // 0-based, so March is 2
        if (isBeforeMarch31) {
            // 1月1日から3月30日まで
            for (let i = 0; i < 4; i++) {
                options.push(`${currentYear + i}卒`);
            }
        } else {
            // 3月31日から12月31日まで
            for (let i = 1; i <= 4; i++) {
                options.push(`${currentYear + i}卒`);
            }
        }
        return options;
    };

    const graduationYearOptions = getGraduationYearOptions();

    if (!isOpen) return null;

    const handleSubmit = () => {
        const newErrors: typeof errors = {};

        if (!firstName.trim()) newErrors.firstName = '必須項目です';
        if (!lastName.trim()) newErrors.lastName = '必須項目です';

        if (!katakanaRegex.test(lastNameKana.trim())) {
            newErrors.lastNameKana = 'カタカナのみ';
        } else if (!lastNameKana.trim()) {
            newErrors.lastNameKana = '必須項目です';
        }

        if (!katakanaRegex.test(firstNameKana.trim())) {
            newErrors.firstNameKana = 'カタカナのみ';
        } else if (!firstNameKana.trim()) {
            newErrors.firstNameKana = '必須項目です';
        }

        if (!schoolName.trim()) newErrors.schoolName = '必須項目です';
        if (!department.trim()) newErrors.department = '必須項目です';
        if (!laboratory.trim()) newErrors.laboratory = '必須項目です';
        if (!graduationYear.trim()) newErrors.graduationYear = '必須項目です';
        if (desiredJobTypes.length === 0) newErrors.desiredJobTypes = '必須項目です';

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        onSubmit(
            firstName,
            lastName,
            firstNameKana,
            lastNameKana,
            schoolName,
            department,
            laboratory,
            graduationYear,
            desiredJobTypes,
            skills
        );
        onClose();
    };

    const handleLogout = () => {
        fetch('http://localhost:8080/auth/logout', {
            method: 'POST',
            credentials: 'include',
        })
            .then((response) => {
                if (response.ok) {
                    // ログアウト成功時の処理
                    router.push('/auth'); // ログインページにリダイレクト
                } else {
                    throw new Error('Failed to logout');
                }
            })
            .catch((error) => {
                console.error('Error during logout:', error);
            });
    };

    const toggleDropdown = () => {
        setIsDropdownOpen(!isDropdownOpen);
    };

    const handleJobTypeChange = (jobType: string) => {
        if (desiredJobTypes.includes(jobType)) {
            setDesiredJobTypes(desiredJobTypes.filter((type) => type !== jobType));
        } else {
            setDesiredJobTypes([...desiredJobTypes, jobType]);
        }
    };

    return (
        <div className="fixed inset-0 flex items-center justify-center bg-orange-500 bg-opacity-50 z-50">
            {/* モーダル全体を包む relative なコンテナを追加 */}
            <div className="relative">
                <div className="bg-white p-8 rounded-lg shadow-lg max-w-screen-2xl w-full overflow-y-auto max-h-[80vh]">
                    <h2 className="text-xl font-bold mb-6">基本情報の登録</h2>
                    <p className="text-gray-700 mb-6">
                        サービスをご利用いただくにあたって、以下の情報を入力してください。
                    </p>
                    {/* 名前 必須のラベル */}
                    <div className="mb-6">
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            名前
                            <span className="bg-orange-500 text-white text-xs font-semibold ml-2 px-2 py-1 rounded">
                                必須
                            </span>
                        </label>
                        <div className="flex space-x-4">
                            <div className="w-1/2">
                                <input
                                    type="text"
                                    placeholder="姓"
                                    value={lastName}
                                    onChange={(e) => setLastName(e.target.value)}
                                    className={`w-full p-2 border ${errors.lastName ? 'border-red-500' : 'border-gray-300'} rounded`}
                                />
                                {errors.lastName && <p className="text-red-500 text-sm mt-1">{errors.lastName}</p>}
                            </div>
                            <div className="w-1/2">
                                <input
                                    type="text"
                                    placeholder="名"
                                    value={firstName}
                                    onChange={(e) => setFirstName(e.target.value)}
                                    className={`w-full p-2 border ${errors.firstName ? 'border-red-500' : 'border-gray-300'} rounded`}
                                />
                                {errors.firstName && <p className="text-red-500 text-sm mt-1">{errors.firstName}</p>}
                            </div>
                        </div>
                    </div>
                    {/* フリガナ 必須のラベル */}
                    <div
                        className="mb-6"
                    >
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            フリガナ
                            <span className="bg-orange-500 text-white text-xs font-semibold ml-2 px-2 py-1 rounded">
                                必須
                            </span>
                        </label>
                        <div className="flex space-x-4">
                            <div className="w-1/2">
                                <input
                                    type="text"
                                    placeholder="セイ"
                                    value={lastNameKana}
                                    onChange={(e) => setLastNameKana(e.target.value)}
                                    className={`w-full p-2 border ${errors.lastNameKana ? 'border-red-500' : 'border-gray-300'} rounded`}
                                />
                                {errors.lastNameKana && <p className="text-red-500 text-sm mt-1">{errors.lastNameKana}</p>}
                            </div>
                            <div className="w-1/2">
                                <input
                                    type="text"
                                    placeholder="メイ"
                                    value={firstNameKana}
                                    onChange={(e) => setFirstNameKana(e.target.value)}
                                    className={`w-full p-2 border ${errors.firstNameKana ? 'border-red-500' : 'border-gray-300'} rounded`}
                                />
                                {errors.firstNameKana && <p className="text-red-500 text-sm mt-1">{errors.firstNameKana}</p>}
                            </div>
                        </div>
                    </div>
                    {/* 学校 必須のラベル */}
                    <div className="mb-6">
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            学校
                            <span className="bg-orange-500 text-white text-xs font-semibold ml-2 px-2 py-1 rounded">
                                必須
                            </span>
                        </label>
                        <input
                            type="text"
                            placeholder="例) ○○大学 ××専門学校"
                            value={schoolName}
                            onChange={(e) => setSchoolName(e.target.value)}
                            className={`w-full p-2 border ${errors.schoolName ? 'border-red-500' : 'border-gray-300'} rounded`}
                        />
                        {errors.schoolName && <p className="text-red-500 text-sm mt-1">{errors.schoolName}</p>}
                    </div>
                    {/* 学部・学科・コース 必須のラベル */}
                    <div className="mb-6">
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            学部・学科・コース
                            <span className="bg-orange-500 text-white text-xs font-semibold ml-2 px-2 py-1 rounded">
                                必須
                            </span>
                        </label>
                        <input
                            type="text"
                            placeholder="例) ○○学科 ××プログラミングコース"
                            value={department}
                            onChange={(e) => setDepartment(e.target.value)}
                            className={`w-full p-2 border ${errors.department ? 'border-red-500' : 'border-gray-300'} rounded`}
                        />
                        {errors.department && <p className="text-red-500 text-sm mt-1">{errors.department}</p>}
                    </div>
                    {/* 研究室 必須 */}
                    <div className="mb-6">
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            研究室
                            <span className="bg-orange-500 text-white text-xs font-semibold ml-2 px-2 py-1 rounded">
                                必須
                            </span>
                        </label>
                        <input
                            type="text"
                            placeholder="プログラミング研究室、田中研究室"
                            value={laboratory}
                            onChange={(e) => setLaboratory(e.target.value)}
                            className={`w-full p-2 border ${errors.laboratory ? 'border-red-500' : 'border-gray-300'} rounded`}
                        />
                        {errors.laboratory && <p className="text-red-500 text-sm mt-1">{errors.laboratory}</p>}
                    </div>
                    {/* 卒業予定年 必須 */}
                    <div className="mb-6">
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            卒業予定年
                            <span className="bg-orange-500 text-white text-xs font-semibold ml-2 px-2 py-1 rounded">
                                必須
                            </span>
                        </label>
                        <select
                            value={graduationYear}
                            onChange={(e) => setGraduationYear(e.target.value)}
                            className={`w-full p-2 border ${errors.graduationYear ? 'border-red-500' : 'border-gray-300'} rounded`}
                        >
                            <option value="">卒業予定年を選択</option>
                            {/* 卒業年度のオプションを動的に生成 */}
                            {graduationYearOptions.map((year) => (
                                <option key={year} value={year}>{year}</option>
                            ))}
                        </select>
                        {errors.graduationYear && <p className="text-red-500 text-sm mt-1">{errors.graduationYear}</p>}
                    </div>
                    {/* 希望職種 必須 */}
                    <div
                        className="mb-6 relative"
                        ref={dropdownRef}
                        onClick={toggleDropdown}
                    >
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            希望職種 （複数可）
                            <span className="bg-orange-500 text-white text-xs font-semibold ml-2 px-2 py-1 rounded">
                                必須
                            </span>
                        </label>
                        <div
                            className={`w-full p-2 border rounded cursor-pointer overflow-hidden text-ellipsis whitespace-nowrap ${errors.desiredJobTypes ? 'border-red-500' : 'border-gray-300'
                                }`}
                            title={desiredJobTypes.join(', ')}
                        >
                            {desiredJobTypes.length > 0
                                ? desiredJobTypes.join(', ')
                                : '希望職種を選択'}
                        </div>
                        {isDropdownOpen && (
                            <div
                                className="absolute top-full left-0 w-full bg-white border border-gray-300 rounded shadow-lg mt-2 max-h-60 overflow-y-auto z-50"
                                onClick={(e) => e.stopPropagation()}
                            >
                                {jobTypeOptions.map((jobType) => (
                                    <div
                                        key={jobType}
                                        className={`p-2 cursor-pointer hover:bg-gray-200 flex items-center`}
                                        onClick={() => handleJobTypeChange(jobType)}
                                    >
                                        <input
                                            type="checkbox"
                                            className="form-checkbox mr-2"
                                            checked={desiredJobTypes.includes(jobType)}
                                            readOnly
                                        />
                                        <span>{jobType}</span>
                                    </div>
                                ))}
                            </div>
                        )}
                        {/* エラーメッセージの表示 */}
                        {errors.desiredJobTypes && (
                            <p className="text-red-500 text-sm mt-1">{errors.desiredJobTypes}</p>
                        )}
                    </div>
                    {/* スキル 任意 */}
                    <div
                        className="mb-6 relative"
                        ref={skillContainerRef}
                        onClick={() => setSuggestedSkills([])}
                    >
                        <label className="text-gray-700 font-bold inline-block mb-2">
                            スキル （複数可）
                            <span className="bg-gray-500 text-white text-xs font-semibold ml-2 px-2 py-1 rounded">
                                任意
                            </span>
                        </label>
                        <div className="border border-gray-300 rounded p-2 relative">
                            {/* Display selected skills as tags */}
                            <div className="flex flex-wrap gap-2 mb-2">
                                {skills.map((skill) => (
                                    <span key={skill} className="bg-gray-200 px-3 py-1 rounded-full text-sm flex items-center">
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
                            <input
                                type="text"
                                placeholder="スキルを選択・追加"
                                value={skillInput}
                                onChange={(e) => setSkillInput(e.target.value)}
                                onKeyDown={handleSkillInputKeyDown}
                                className="w-full p-2 border border-gray-300 rounded"
                            />
                            {skillInput && (
                                <div
                                    className="absolute left-0 z-50 bg-white border border-gray-300 rounded mt-2 w-full max-h-40 overflow-y-auto"
                                    onClick={(e) => e.stopPropagation()}
                                >
                                    {suggestedSkills.length > 0 ? (
                                        suggestedSkills.map((skill) => (
                                            <div
                                                key={skill}
                                                className="p-2 hover:bg-gray-200 cursor-pointer"
                                                onClick={() => handleSkillSelect(skill)}
                                            >
                                                {skill}
                                            </div>
                                        ))
                                    ) : (
                                        <div
                                            className="p-2 hover:bg-gray-200 cursor-pointer"
                                            onClick={() => handleSkillSelect(skillInput.trim())}
                                        >
                                            Enterで追加
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                    {/* 注意書きの追加 */}
                    <p className="text-xs text-gray-500 mt-4">
                        ※入力内容は後で変更することができます。<br />
                        ※登録情報は正確かつ最新の内容となるよう、適宜アップデートをお願いします。
                    </p>
                    <button
                        className="w-full py-2 bg-orange-500 text-white rounded-md hover:bg-orange-600 mt-4"
                        onClick={handleSubmit}
                    >
                        登録する
                    </button>
                </div>
                {/* ログアウトボタンをモーダルボディの外側に配置 */}
                <button
                    className="absolute -bottom-10 right-0 text-white px-4 py-2 rounded-md hover:text-gray-300"
                    onClick={handleLogout}
                >
                    ログアウト
                </button>
            </div>
        </div>
    );
};

export default MnimumUserInfoInputModal;
