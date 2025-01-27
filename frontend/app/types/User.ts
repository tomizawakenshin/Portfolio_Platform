export interface User {
    ID: number;
    FirstName: string;
    LastName: string;
    FirstNameKana: string;
    LastNameKana: string
    Email: string;
    profilePictureURL?: string;
    GraduationYear: number;
}