import styles from "../../../components/editor/TiptapEditor.module.css"
import type { CreateOrUpdateProblemRequest } from "../../../types/request"

interface ProblemPreviewProps {
  formData: CreateOrUpdateProblemRequest
}

export default function ProblemPreview({ formData }: ProblemPreviewProps) {
  const difficultyColors = {
    easy: "text-green-600 bg-green-100",
    medium: "text-yellow-600 bg-yellow-100",
    hard: "text-red-600 bg-red-100",
  }

  return (
    <div className="space-y-6 bg-white rounded-lg p-6 shadow-sm">
      <div className="border-b pb-4">
        <div className="flex items-center justify-between mb-2">
          <h1 className="text-3xl font-bold">
            {formData.title || "Untitled Problem"}
          </h1>
          <span className={`px-3 py-1 rounded-full text-sm font-medium ${difficultyColors[formData.difficulty]}`}>
            {formData.difficulty.charAt(0).toUpperCase() + formData.difficulty.slice(1)}
          </span>
        </div>
        <p className="text-gray-500 text-sm">/{formData.slug || "problem-slug"}</p>
      </div>

      {/* Description */}
      {formData.description && (
        <div>
          <h2 className="text-xl font-semibold mb-3">Description</h2>
          <div 
            className={styles.editor}
            dangerouslySetInnerHTML={{ __html: formData.description }}
          />
        </div>
      )}

      {/* Input Format */}
      {formData.input_format && (
        <div>
          <h2 className="text-xl font-semibold mb-3">Input Format</h2>
          <div 
            className={styles.editor}
            dangerouslySetInnerHTML={{ __html: formData.input_format }}
          />
        </div>
      )}

      {/* Output Format */}
      {formData.output_format && (
        <div>
          <h2 className="text-xl font-semibold mb-3">Output Format</h2>
          <div 
            className={styles.editor}
            dangerouslySetInnerHTML={{ __html: formData.output_format }}
          />
        </div>
      )}

      {/* Constraints */}
      {formData.constraints && (
        <div>
          <h2 className="text-xl font-semibold mb-3">Constraints</h2>
          <div 
            className={styles.editor}
            dangerouslySetInnerHTML={{ __html: formData.constraints }}
          />
        </div>
      )}

      {/* Limits */}
      <div className="bg-gray-50 rounded-lg p-4">
        <h2 className="text-lg font-semibold mb-3">Limits</h2>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <span className="text-gray-600">Time Limit:</span>
            <span className="ml-2 font-medium">{formData.time_limit} ms</span>
          </div>
          <div>
            <span className="text-gray-600">Memory Limit:</span>
            <span className="ml-2 font-medium">{formData.memory_limit} MB</span>
          </div>
        </div>
      </div>

      {/* Status */}
      <div className="flex gap-4 text-sm">
        <span className={`px-3 py-1 rounded ${formData.status === "published" ? "bg-blue-100 text-blue-700" : "bg-gray-100 text-gray-700"}`}>
          {formData.status.charAt(0).toUpperCase() + formData.status.slice(1)}
        </span>
        {formData.is_active && (
          <span className="px-3 py-1 rounded bg-green-100 text-green-700">
            Active
          </span>
        )}
      </div>
    </div>
  )
}
