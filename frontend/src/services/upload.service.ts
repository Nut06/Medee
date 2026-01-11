import { api } from "@/lib/api";
import axios from "axios";

interface SignedURLResponse {
  uploadUrl: string;
  path: string;
  publicUrl: string;
}

export const uploadService = {
  async getSignedUploadUrl(filename: string): Promise<SignedURLResponse> {
    const { data } = await api.post("upload/signed-url", { filename });
    return data as SignedURLResponse;
  },

  async uploadImage(file: File): Promise<string> {
    const { uploadUrl, publicUrl } = await this.getSignedUploadUrl(file.name);
    await axios.put(uploadUrl, file, {
      headers: {
        "Content-Type": file.type,
      },
    });
    console.log("Upload Image success");
    return publicUrl;
  },

  async uploadImages(files: File[]): Promise<string[]> {
    return Promise.all(files.map((file) => this.uploadImage(file)));
  },
};
