import { apiRequest } from "./api";

export type Document = {
  docID: string;
  title: string;
  content: string;
  ownerID: string;
  version: number;
  createdAt: string;
  updatedAt: string;
  deletedAt: {
    Time: string;
    Valid: boolean;
};
  updatedBy: string;
};

export type DocumentVersion = {
  id: string;
  docID: string;
  versionID: number;
  title: string;
  content: string;
  modifiedByUserID: string;
  modifiedByUsername: string;
  createdAt: string;
};

type ApiResponse<T> = {
  message: string;
  data: T;
};

type CreateDocumentResponse = {
  docID: string;
};

export async function getMyDocuments(token: string): Promise<Document[]> {
  const response = await apiRequest<ApiResponse<Document[]>>(
    "/docs/all",
    { method: "GET" },
    token,
  );

  return response.data;
}

export async function createDocument( title: string, content: string, token: string ): Promise<CreateDocumentResponse> {
  const response = await apiRequest<ApiResponse<CreateDocumentResponse>>(
    "/docs",
    {
      method: "POST",
      body: JSON.stringify({ title, content }),
    },
    token,
  );

  return response.data;
}

export async function getDocument( docID: string, token: string ): Promise<Document> {
  const response = await apiRequest<ApiResponse<Document>>(`/docs/${docID}`,
    { method: "GET" },
    token,
  );

  if (!response.data) {
    throw new Error("Document data was not returned by the server.");
  }

  return response.data;
}

export async function updateDocument( docID: string, title: string, content: string, token: string ): Promise<void> {
  await apiRequest<ApiResponse<unknown>>(`/docs/${docID}`,
    {
      method: "PUT",
      body: JSON.stringify({ title, content }),
    },
    token,
  );
}

export async function saveVersion(docID: string, token: string): Promise<void> {
  await apiRequest<ApiResponse<unknown>>(`/docs/${docID}/versions/save`,
    {
      method: "POST",
    },
    token,
  );
}

export async function getVersions(docID: string, token: string): Promise<DocumentVersion[]> {
  const response = await apiRequest<ApiResponse<DocumentVersion[]>>(`/docs/${docID}/versions`,
    { method: "GET" },
    token,
  );

  return response.data;
}

export async function getVersion(docID: string, versionID: number, token: string): Promise<DocumentVersion> {
  const response = await apiRequest<ApiResponse<DocumentVersion>>(`/docs/${docID}/versions/${versionID}`,
    { method: "GET" },
    token,
  );

  if (!response.data) {
    throw new Error("Version data was not returned by the server.");
  }

  return response.data;
}

export async function restoreVersion(docID: string, versionID: number, token: string): Promise<void> {
  await apiRequest<ApiResponse<unknown>>(`/docs/${docID}/versions/${versionID}/restore`,
    { method: "POST" },
    token,
  );
}