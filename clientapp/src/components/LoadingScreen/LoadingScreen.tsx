import styles from "./LoadingScreen.module.css";

function LoadingScreen() {
    return (
      <div className={styles.fullScreen}>
          <p className={styles.centeredText}>Securing...</p>
      </div>
    );
}

export default LoadingScreen;
