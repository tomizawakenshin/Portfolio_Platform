'use client';

import { useState } from "react";
import Header from "../components/Header";
import MainContent from "../components/MainContent";
import SignUpModal from "../components/SignUpModal";
import LoginModal from "../components/LoginModal";
import EmailSignUpModal from "../components/EmailSignUpModal";
import SignUpCompleteModal from "../components/SignUpCompleteModal";
import ForgotPasswordModal from "../components/ForgotPasswordModal";

export default function Home() {
  const [isModalOpen, setIsModalOpen] = useState(false);       // サインアップモーダル
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false); // ログインモーダル
  const [isEmailSignUpModalOpen, setIsEmailSignUpModalOpen] = useState(false); // 新しいメールサインアップモーダル
  const [isCompleteModalOpen, setIsCompleteModalOpen] = useState(false);
  const [userEmail, setUserEmail] = useState("");
  const [isForgotPasswordModalOpen, setIsForgotPasswordModalOpen] = useState(false);

  const openSignUpModal = () => {
    setIsModalOpen(true);
  };

  const closeSignUpModal = () => {
    setIsModalOpen(false);
  };

  const openLoginModal = () => {
    setIsLoginModalOpen(true);
  };

  const closeLoginModal = () => {
    setIsLoginModalOpen(false);
  };

  // SignUpModalでログインボタンをクリックしたときの処理
  const handleSwitchToLogin = () => {
    closeSignUpModal();       // サインアップモーダルを閉じる
    openLoginModal();   // ログインモーダルを開く
  };

  const handleSwitchToSignUp = () => {
    closeLoginModal();
    openSignUpModal();
  }

  const openEmailSignUpModal = () => {
    setIsEmailSignUpModalOpen(true);
    closeSignUpModal();
  };

  const closeEmailSignUpModal = () => {
    setIsEmailSignUpModalOpen(false);
  };

  const handleSignUpComplete = (email: string) => {
    setUserEmail(email);
    setIsCompleteModalOpen(true);
  };

  const closeCompleteModal = () => {
    setIsCompleteModalOpen(false);
  };

  // ForgotPasswordModalを開く関数を定義
  const openForgotPasswordModal = () => {
    setIsLoginModalOpen(false); // LoginModalを閉じる
    setIsForgotPasswordModalOpen(true); // ForgotPasswordModalを開く
  };

  // ForgotPasswordModalを閉じる関数を定義
  const closeForgotPasswordModal = () => {
    setIsForgotPasswordModalOpen(false);
  };

  return (
    <div className="font-sans">
      {/* Header */}
      <Header />

      {/* Main content */}
      <MainContent
        onFreeStartClick={openSignUpModal}
        onLoginClick={openLoginModal}
      />

      {/* SignUp Modal */}
      <SignUpModal
        isOpen={isModalOpen}
        onClose={closeSignUpModal}
        onLoginClick={handleSwitchToLogin} // ログインボタンがクリックされたときのハンドラ
        onEmailSignUpClick={openEmailSignUpModal} // メールサインアップボタンを追加
      />

      {/* Login Modal */}
      <LoginModal
        isOpen={isLoginModalOpen}
        onClose={closeLoginModal}
        onSignUpClick={handleSwitchToSignUp}
        onForgotPasswordClick={openForgotPasswordModal}
      />

      <EmailSignUpModal
        isOpen={isEmailSignUpModalOpen}
        onClose={closeEmailSignUpModal}
        onComplete={handleSignUpComplete}
      />

      <SignUpCompleteModal
        isOpen={isCompleteModalOpen}
        onClose={closeCompleteModal}
        email={userEmail}
      />

      <ForgotPasswordModal
        isOpen={isForgotPasswordModalOpen}
        onClose={closeForgotPasswordModal}
      />
    </div>
  );
}
