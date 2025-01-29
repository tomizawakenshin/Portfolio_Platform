export interface User {
    ID: number;
    FirstName: string;
    LastName: string;
    FirstNameKana: string;
    LastNameKana: string
    Email: string;
    SchoolName: string;
    Department: string;
    Laboratory: string;
    GraduationYear: number;
    DesiredJobTypes: string[];
    Skills: string[];
    SelfIntroduction: string;
    ProfileImageURL: string;
}