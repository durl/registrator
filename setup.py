#!/usr/bin/env python

from distutils.core import setup

__version__ = '0.1.0'

with open('README.rst') as readme:
    readme = readme.read()

setup(
    name='smartmeter-analyze',
    version=__version__,
    description='',
    long_description=readme,
    author='David Url',
    author_email='david@x00.at',
    url='https://github.com/durl/smartmeter-analyze',
    py_modules=[
        'registrator',
    ],
    requires=[
    ],
    license='License :: OSI Approved :: BSD License',
    keywords='',
    classifiers=[
        'Environment :: No Input/Output (Daemon)',
        'Intended Audience :: Developers',
        'Intended Audience :: System Administrators',
        'License :: OSI Approved :: BSD License',
        'Operating System :: POSIX :: Linux',
        'Programming Language :: Python :: 2.7',
        'Topic :: System :: Networking :: Monitoring',
        'Topic :: System :: Systems Administration',
    ],
)
