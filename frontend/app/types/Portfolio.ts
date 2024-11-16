// types/Post.ts

import { User } from './User';

export interface Image {
    ID: number;
    URL: string;
    PostID: number;
}

export interface Portfolio {
    ID: number;
    Title: string;
    Description: string;
    GitHubLink: string;
    ProductLink: string;
    Skills: string[];
    Images: Image[];
    UserID: number;
    User: User;
}
