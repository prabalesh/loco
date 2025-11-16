import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Code, Trophy, Users, ArrowRight } from 'lucide-react'
import { Button } from '../shared/components/ui/Button'
import { Card } from '../shared/components/ui/Card'
import { useAuth } from '../shared/hooks/useAuth'
import { ROUTES } from '../shared/constants/routes'
import { CONFIG } from '../shared/constants/config'

const features = [
  {
    icon: Code,
    title: '500+ Problems',
    description: 'From easy to hard, covering all algorithms and data structures',
  },
  {
    icon: Trophy,
    title: 'Weekly Contests',
    description: 'Compete with developers worldwide and climb the leaderboard',
  },
  {
    icon: Users,
    title: 'Active Community',
    description: 'Discuss solutions and learn from others',
  },
]

const stats = [
  { label: 'Active Users', value: '10K+' },
  { label: 'Problems Solved', value: '500K+' },
  { label: 'Code Submissions', value: '2M+' },
]

export const HomePage = () => {
  const { isAuthenticated } = useAuth()

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="bg-gradient-to-br from-blue-600 to-indigo-700 text-white py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="text-center"
          >
            <h1 className="text-5xl md:text-6xl font-bold mb-6">
              Master Coding Through
              <br />
              <span className="text-blue-200">Practice & Competition</span>
            </h1>
            <p className="text-xl md:text-2xl mb-8 text-blue-100 max-w-3xl mx-auto">
              {CONFIG.APP_DESCRIPTION}. Join thousands of developers improving their skills.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              {isAuthenticated ? (
                <Link to={ROUTES.PROBLEMS}>
                  <Button variant="primary" size="lg" className="bg-white !text-blue-600 hover:!bg-gray-100 shadow-lg">
                    Start Solving
                    <ArrowRight className="ml-2 h-5 w-5" />
                  </Button>
                </Link>
              ) : (
                <>
                  <Link to={ROUTES.REGISTER}>
                    <Button variant="primary" size="lg" className="bg-white !text-blue-600 hover:!bg-gray-100 shadow-lg">
                      Get Started Free
                      <ArrowRight className="ml-2 h-5 w-5" />
                    </Button>
                  </Link>
                  <Link to={ROUTES.PROBLEMS}>
                    <Button variant="outline" size="lg" className="border-white text-white hover:bg-white hover:text-blue-600">
                      Explore Problems
                    </Button>
                  </Link>
                </>
              )}
            </div>
          </motion.div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="bg-white py-12 border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {stats.map((stat, index) => (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                className="text-center"
              >
                <div className="text-4xl font-bold text-blue-600 mb-2">
                  {stat.value}
                </div>
                <div className="text-gray-600 font-medium">{stat.label}</div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.6 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              Why Choose {CONFIG.APP_NAME}?
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Everything you need to ace coding interviews and become a better developer
            </p>
          </motion.div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <motion.div
                key={feature.title}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                viewport={{ once: true }}
              >
                <Card className="p-8 h-full hover:shadow-lg transition-shadow">
                  <feature.icon className="h-12 w-12 text-blue-600 mb-4" />
                  <h3 className="text-xl font-bold text-gray-900 mb-2">
                    {feature.title}
                  </h3>
                  <p className="text-gray-600">{feature.description}</p>
                </Card>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.6 }}
            viewport={{ once: true }}
            className="text-center mb-16"
          >
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              How It Works
            </h2>
            <p className="text-xl text-gray-600">
              Start solving problems in just 3 simple steps
            </p>
          </motion.div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
            {[
              {
                step: '1',
                title: 'Choose a Problem',
                description: 'Browse our curated list of coding challenges',
              },
              {
                step: '2',
                title: 'Write Your Solution',
                description: 'Use our online editor with syntax highlighting',
              },
              {
                step: '3',
                title: 'Submit & Learn',
                description: 'Get instant feedback and compare solutions',
              },
            ].map((item, index) => (
              <motion.div
                key={item.step}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.5, delay: index * 0.2 }}
                viewport={{ once: true }}
                className="text-center"
              >
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-blue-600 text-white text-2xl font-bold mb-4">
                  {item.step}
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-2">
                  {item.title}
                </h3>
                <p className="text-gray-600">{item.description}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-gradient-to-r from-blue-600 to-indigo-700 text-white py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            viewport={{ once: true }}
          >
            <h2 className="text-4xl font-bold mb-4">
              Ready to Start Coding?
            </h2>
            <p className="text-xl mb-8 text-blue-100">
              Join thousands of developers improving their skills today
            </p>
            {!isAuthenticated && (
              <Link to={ROUTES.REGISTER}>
                <Button variant="primary" size="lg" className="bg-white !text-blue-600 hover:bg-gray-100">
                  Sign Up - It's Free
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
              </Link>
            )}
          </motion.div>
        </div>
      </section>
    </div>
  )
}
